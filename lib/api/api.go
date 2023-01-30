package api

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"math/rand"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/pool"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/code-to-go/safepool.lib/sql"
	"gopkg.in/yaml.v3"
)

var Apps = []string{
	"chat",
	"library",
}

var Self security.Identity

//go:embed sqlite.sql
var sqlliteDDL string

// SetDbPath set the full path where the DB will be created. Useful on Android/iOS platforms
func SetDbPath(dbPath string) {
	sql.DbPath = dbPath
}

// SetDbName set the name of the DB file. Useful for test purpose
func SetDbName(dbName string) {
	sql.DbPath = dbName
}

func Start(dbPath string) error {
	sql.InitDDL = sqlliteDDL

	err := sql.OpenDB(dbPath)
	if core.IsErr(err, "cannot open DB: %v") {
		return err
	}

	s, _, _, ok := sqlGetConfig("", "SELF")
	if ok {
		Self, err = security.IdentityFromBase64(s)
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		name := fmt.Sprintf("%s%d", names[r.Intn(len(names))], r.Intn(100))
		Self, err = security.NewIdentity(name)
		if err == nil {
			s, err = Self.Base64()
			if err == nil {
				err = sqlSetConfig("", "SELF", s, 0, nil)
			}
			if core.IsErr(err, "çannot save identity to db: %v") {
				panic("cannot save identity in db")
			}

			err = security.SetIdentity(Self)
			if core.IsErr(err, "çannot save identity to db: %v") {
				panic("cannot save identity in db")
			}

			err = security.Trust(Self, true)
			if core.IsErr(err, "çannot set trust of '%s' on db: %v", Self.Nick) {
				panic("cannot trust Self in db")
			}

		}

	}
	return err
}

func SetNick(nick string) error {
	Self.Nick = nick
	s, err := Self.Base64()
	if core.IsErr(err, "cannot serialize self to db: %v") {
		return err
	}
	err = sqlSetConfig("", "SELF", s, 0, nil)
	core.IsErr(err, "cannot save nick to db: %v")
	return err
}

func CreatePool(c pool.Config, apps []string) error {
	err := pool.Define(c)
	if core.IsErr(err, "cannot define pool %v: %s", c.Name) {
		return err
	}

	p, err := pool.Create(Self, c.Name, apps)
	if core.IsErr(err, "cannot create pool %v: %s", c.Name) {
		return err
	}
	p.Close()
	return nil
}

func AddPool(token string) error {
	var c pool.Config
	bs, err := base64.StdEncoding.DecodeString(token)
	if core.IsErr(err, "cannot decousterde token: %v") {
		return err
	}
	err = yaml.Unmarshal(bs, &c)
	if core.IsErr(err, "cannot unmarshal config from token: %v") {
		return err
	}

	if c.Name == "" || (len(c.Public)+len(c.Private)) == 0 {
		core.IsErr(pool.ErrInvalidToken, "invalid config '%s': %v", string(bs))
		return pool.ErrInvalidToken
	} else {
		core.Info("valid token for pool '%s'", c.Name)
	}

	err = pool.Define(c)
	if core.IsErr(err, "cannot define pool '%s': %v", c.Name) {
		return err
	}
	p, err := pool.Open(Self, c.Name)
	if core.IsErr(err, "cannot open pool '%s': %v", c.Name) {
		return err
	}
	p.Close()

	return nil
}
