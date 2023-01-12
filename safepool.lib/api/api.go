package api

import (
	_ "embed"
	"fmt"
	"math/rand"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/code-to-go/safepool.lib/sql"
)

var Self security.Identity

//go:embed sqlite.sql
var sqlliteDDL string

func Start() {
	sql.InitDDL = sqlliteDDL

	err := sql.OpenDB()
	if err != nil {
		panic("cannot open DB")
	}

	s, _, _, ok := sqlGetConfig("", "SELF")
	if ok {
		Self, err = security.IdentityFromBase64(s)
	} else {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		name := fmt.Sprintf("%s #%d", names[r.Intn(len(names))], r.Intn(100))
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
	if err != nil {
		panic("corrupted identity in DB")
	}

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
