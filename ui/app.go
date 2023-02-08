package main

import (
	"context"
	"fmt"
	"path/filepath"
	"strconv"

	"github.com/adrg/xdg"
	"github.com/code-to-go/safe/safepool/apps/chat"
	"github.com/code-to-go/safe/safepool/apps/library"
	"github.com/code-to-go/safe/safepool/core"
	"github.com/code-to-go/safe/safepool/pool"
	"github.com/code-to-go/safe/safepool/security"
	"github.com/code-to-go/safe/safepool/transport"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// App struct
type App struct {
	ctx context.Context
}

// NewApp creates a new App application struct
func NewApp() *App {
	dbPath := filepath.Join(xdg.ConfigHome, "safepool.test.db")
	core.FatalIf(safepool.Start(dbPath), "cannot initialize pool: %v")

	for _, n := range pool.List() {
		p, err := pool.Open(safepool.Self, n)
		if err == nil {
			pools[n] = p
		}
	}
	return &App{}
}

var pools = map[string]*pool.Pool{}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx

}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) GetNick() string {
	return safepool.Self.Nick
}

func (a *App) GetPoolList() []string {
	var names []string
	for n := range pools {
		names = append(names, n)
	}
	return names
}

func (a *App) GetConfigTemplate() (string, error) {
	c := pool.Config{
		Name:   "pool name",
		Public: transport.SampleConfig,
	}

	data, err := yaml.Marshal(c)
	return string(data), err
}

func (a *App) CreatePool(config string) (string, error) {
	var c pool.Config
	err := yaml.Unmarshal([]byte(config), &c)
	if err != nil {
		return "", err
	}

	if (len(c.Public)+len(c.Private)) == 0 || c.Name == "" {
		return "", pool.ErrInvalidConfig
	} else {
		core.Info("valid config for pool '%s'", c.Name)
	}
	err = pool.Define(c)
	if core.IsErr(err, "cannot define pool '%s': %v", c.Name) {
		return "", err
	}

	p, err := pool.Create(safepool.Self, c.Name, safepool.Apps)
	if core.IsErr(err, "cannot create pool '%s': %v", c.Name) {
		return "", err
	}
	pools[p.Name] = p

	token, err := pool.EncodeToken(pool.Token{
		Config: c,
		Host:   safepool.Self,
	}, "")
	core.IsErr(err, "cannot encode universal token: %v")

	return token, err
}

func (a *App) AddPool(token string) error {
	_, err := safepool.AddPool(token)
	return err
}

func (a *App) GetMessages(poolName string, afterIdS string, beforeIdS string, limit int) ([]chat.Message, error) {
	beforeId, err := strconv.ParseUint(beforeIdS, 10, 64)
	if core.IsErr(err, "invalid beforeId parameter '%s': %v", beforeIdS) {
		return nil, err
	}
	afterId, err := strconv.ParseUint(afterIdS, 10, 64)
	if core.IsErr(err, "invalid afterId parameter '%s': %v", afterIdS) {
		return nil, err
	}

	if p, ok := pools[poolName]; ok {
		c := chat.Get(p, "chat")
		return c.GetMessages(afterId, beforeId, limit)
	}
	return nil, fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) PostMessage(poolName string, text string, contentType string, binary []byte) (string, error) {
	if p, ok := pools[poolName]; ok {
		c := chat.Get(p, "chat")
		id, err := c.SendMessage(contentType, text, binary)
		if core.IsErr(err, "cannot post chat message: %v") {
			return "", err
		}
		return strconv.FormatUint(id, 10), nil
	}
	return "", fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) GetToken(poolName string, guestId string) (string, error) {
	if p, ok := pools[poolName]; ok {
		c, err := pool.GetConfig(poolName)
		if core.IsErr(err, "cannot get pool config: %v") {
			return "", err
		}

		t := pool.Token{
			Config: c,
			Host:   safepool.Self,
		}

		if guestId != "" {
			err = p.SetAccess(guestId, pool.Active)
			if core.IsErr(err, "cannot set access for id '%s' in pool '%s': %v", guestId, p.Name) {
				return "", err
			}
		}
		return pool.EncodeToken(t, guestId)
	}
	return "", fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) ListLibrary(poolName string) []library.File {
	return []library.File{
		{
			Id:          1,
			Name:        "test.doc",
			Size:        192922,
			AuthorId:    safepool.Self.Id(),
			ContentType: "application/word",
			// LocalPath:   "~/Documents/pools/safepool.ch/test.doc",
			// Mode:        library.Sync,
		},
		{
			Id:          1,
			Name:        "test2.doc",
			Size:        192922,
			AuthorId:    safepool.Self.Id(),
			ContentType: "application/word",
			// LocalPath:   "~/Documents/pools/safepool.ch/test.doc",
			// Mode:        library.Sync,
		},
	}
}

func (a *App) GetIdentities(poolName string) ([]security.Identity, error) {
	if p, ok := pools[poolName]; ok {
		return p.Identities()
	}
	return nil, fmt.Errorf("invalid pool '%s'", poolName)
}

func (a *App) UpdateIdentity(identity security.Identity) error {
	return security.SetIdentity(identity)
}

func (a *App) GetSelf() security.Identity {
	return safepool.Self.Public()
}

func (a *App) GetSelfId() string {
	return safepool.Self.Id()
}

func (a *App) DecodeToken(token string) (pool.Token, error) {
	t, err := pool.DecodeToken(safepool.Self, token)
	logrus.Infof("%v", t)
	return t, err
}
