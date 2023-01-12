package pool

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/code-to-go/safepool.lib/transport"

	"github.com/godruoyi/go-snowflake"
)

const SafeConfigFile = ".safepool-pool.json"

var ErrNoExchange = errors.New("no Exchange available")
var ErrInvalidSignature = errors.New("signature is invalid")
var ErrNotTrusted = errors.New("the author is not a trusted user")
var ErrNotAuthorized = errors.New("no authorization for this file")
var ErrAlreadyExist = errors.New("pool already exists")
var ErrInvalidToken = errors.New("provided token is invalid: missing name or configs")
var ErrInvalidId = errors.New("provided id not a valid ed25519 public key")
var ErrInvalidConfig = errors.New("provided config is invalid: missing name or configs")

type Consumer interface {
	TimeOffset(s *Pool) time.Time
	Accept(s *Pool, h Head) bool
}

type Pool struct {
	Name      string
	Self      security.Identity
	Consumers []Consumer

	e           transport.Exchanger
	exchangers  []transport.Exchanger
	masterKeyId uint64
	masterKey   []byte
	lastReplica time.Time
	accessHash  []byte
}

type Identity struct {
	security.Identity
	//Since is the keyId used when the identity was added to the Pool access
	Since uint64
	//AddedOn is the timestamp when the identity is stored on the local DB
	AddedOn time.Time
}

type Head struct {
	Id        uint64
	Name      string
	Size      int64
	Hash      []byte
	ModTime   time.Time
	AuthorId  string
	Signature []byte
	Meta      []byte
	Offset    int `json:"-"`
}

const (
	ID_CREATE       = 0x0
	ID_FORCE_CREATE = 0x1
)

var ForceCreation = false
var ReplicaPeriod = time.Hour
var CacheSizeMB = 16

type Config struct {
	Name    string
	Configs []transport.Config
}

func List() []string {
	names, _ := sqlList()
	return names
}

func Create(self security.Identity, name string) (*Pool, error) {
	configs, err := sqlLoad(name)
	if core.IsErr(err, "unknown pool %s: %v", name) {
		return nil, err
	}

	p := &Pool{
		Name:        name,
		Self:        self,
		lastReplica: time.Now(),
	}
	err = p.connectSafe(name, configs)
	if err != nil {
		return nil, err
	}

	p.masterKeyId = snowflake.ID()
	p.masterKey = security.GenerateBytesKey(32)
	err = p.sqlSetKey(p.masterKeyId, p.masterKey)
	if core.IsErr(err, "çannot store master encryption key to db: %v") {
		return nil, err
	}

	access := Access{
		Id:      self.Id(),
		State:   Active,
		ModTime: time.Now(),
	}
	err = p.sqlSetAccess(access)
	if core.IsErr(err, "cannot link identity to pool '%s': %v", p.Name) {
		return nil, err
	}

	if !ForceCreation {
		_, err = p.e.Stat(path.Join(p.Name, ".access"))
		if err == nil {
			return nil, ErrAlreadyExist
		}
	}

	err = p.syncIdentities()
	if core.IsErr(err, "cannot sync own identity: %v") {
		return nil, err
	}

	err = p.exportAccessFile()
	if core.IsErr(err, "cannot export access file for pool '%s': %v", name) {
		return nil, err
	}

	return p, err
}

// Init initialized a domain on the specified exchangers
func Open(self security.Identity, name string) (*Pool, error) {
	configs, err := sqlLoad(name)
	if core.IsErr(err, "unknown pool %s: %v", name) {
		return nil, err
	}
	p := &Pool{
		Name: name,
		Self: self,
	}
	err = p.connectSafe(name, configs)
	if err != nil {
		return nil, err
	}

	err = p.syncIdentities()
	if core.IsErr(err, "cannot sync own identity: %v") {
		return nil, err
	}

	_, err = p.sync(p.e)
	return p, err
}

type AcceptFunc func(head Head)

const All = ""

func (p *Pool) List(offset int) ([]Head, error) {
	hs, err := sqlGetHeads(p.Name, offset)
	if core.IsErr(err, "cannot read Pool heads: %v") {
		return nil, err
	}
	return hs, err
}

func (p *Pool) Send(name string, r io.Reader, meta []byte) (Head, error) {
	id := snowflake.ID()
	n := path.Join(p.Name, fmt.Sprintf("%d.body", id))
	hr, err := p.writeFile(n, r)
	if core.IsErr(err, "cannot post file %s to %s: %v", name, p.e) {
		return Head{}, err
	}

	hash := hr.Hash()
	signature, err := security.Sign(p.Self, hash)
	if core.IsErr(err, "cannot sign file %s.body in %s: %v", name, p.e) {
		return Head{}, err
	}
	h := Head{
		Id:        id,
		Name:      name,
		Size:      hr.Size(),
		Hash:      hash,
		ModTime:   time.Now(),
		AuthorId:  p.Self.Id(),
		Signature: signature,
		Meta:      meta,
	}
	data, err := json.Marshal(h)
	if core.IsErr(err, "cannot marshal header to json: %v") {
		return Head{}, err
	}

	n = path.Join(p.Name, fmt.Sprintf("%d.head", id))
	_, err = p.writeFile(n, bytes.NewBuffer(data))
	core.IsErr(err, "cannot write header %s.head in %s: %v", name, p.e)

	return h, nil
}

func (p *Pool) Receive(id uint64, rang *transport.Range, w io.Writer) error {
	h, err := sqlGetHead(p.Name, id)
	if err != nil {
		headName := path.Join(p.Name, fmt.Sprintf("%d.head", id))
		h, err = p.readHead(headName)
		if core.IsErr(err, "cannot read header '%s': %v") {
			return err
		}
	}

	bodyName := path.Join(p.Name, fmt.Sprintf("%d.body", id))
	cached, err := p.getFromCache(bodyName, rang, w)
	if cached {
		return err
	}
	cw, err := p.cacheWriter(bodyName, w)
	if err == nil {
		defer cw.Close()
		w = cw
	}

	hr, err := p.readFile(bodyName, rang, w)
	if core.IsErr(err, "cannot read body '%s': %v", bodyName) {
		return err
	}
	hash := hr.Hash()
	if !bytes.Equal(hash, h.Hash) {
		return ErrInvalidSignature
	}

	return nil
}

func (p *Pool) Sync() error {
	if !p.e.Touched(p.Name + "/") {
		return nil
	}
	hs, _ := p.List(0)

	heads := map[uint64]Head{}
	for _, h := range hs {
		heads[h.Id] = h
	}

	fs, err := p.e.ReadDir(p.Name, 0)
	if core.IsErr(err, "cannot read content in pool %s/%s", p.e, p.Name) {
		return err
	}
	for _, f := range fs {
		name := f.Name()
		if !strings.HasSuffix(name, ".head") {
			continue
		}

		id, err := strconv.ParseInt(name[0:len(name)-len(".head")], 10, 64)
		if err != nil {
			continue
		}
		if _, found := heads[uint64(id)]; found {
			continue
		}

		h, err := p.readHead(path.Join(p.Name, name))
		if core.IsErr(err, "cannot read file %s from %s: %v", name, p.e) {
			continue
		}
		_ = sqlAddHead(p.Name, h)
		hs = append(hs, h)
	}

	if time.Until(p.lastReplica) > ReplicaPeriod {
		p.replica()
		p.lastReplica = time.Now()
	}
	return nil
}

func (p *Pool) Close() {
	for _, e := range p.exchangers {
		_ = e.Close()
	}
}

func (p *Pool) Delete() error {
	for _, e := range p.exchangers {
		err := e.Delete(p.Name)
		if err != nil {
			return err
		}
	}
	return nil
}

func (p *Pool) Identities() ([]security.Identity, error) {
	identities, _, err := p.sqlGetAccesses(false)
	return identities, err
}

func (p *Pool) SetAccess(userId string, state State) error {
	_, ok, _ := security.GetIdentity(userId)
	if !ok {
		identity, err := security.IdentityFromId(userId)
		if core.IsErr(err, "id '%s' is invalid: %v") {
			return err
		}
		identity.Nick = "❓ Incognito..."
		err = security.SetIdentity(identity)
		if core.IsErr(err, "cannot save identity '%s' to db: %v", identity) {
			return err
		}
	}

	err := p.sqlSetAccess(Access{
		Id:      userId,
		State:   state,
		ModTime: time.Now(),
	})
	if core.IsErr(err, "cannot link identity '%s' to pool '%s': %v", userId, p.Name) {
		return err
	}

	return p.exportAccessFile()
}
