package library

import (
	"bytes"
	"encoding/json"
	"os"
	"path"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	pool "github.com/code-to-go/safepool.lib/pool"
	"github.com/wailsapp/mimetype"
)

type State int

var HashChainMaxLength = 32

var Auto = ""

const (
	StateSync State = iota
	StateIn
	StateOut
	StateAlt
)

type Mode int

const (
	Same Mode = 1 << iota
	Newer
	Older
	Deleted
	Conflict
	Folder
)

// Document includes information about a file stored on the library. Most information refers on the synchronized state with the exchange.
type Document struct {
	Id           uint64    `json:"id"`
	Name         string    `json:"name"`
	Mode         Mode      `json:"mode"`
	ModTime      time.Time `json:"poolTime"`
	Size         uint64    `json:"size"`
	AuthorId     string    `json:"authorId"`
	ContentType  string    `json:"contentType"`
	Hash         []byte    `json:"hash"`
	HashChain    [][]byte  `json:"hashChain"`
	Tags         []string  `json:"tags"`
	LocalModTime time.Time `json:"modTime"`
	LocalPath    string    `json:"localPath"`
	Offset       int       `json:"offset"`
}

type Library struct {
	Pool    *pool.Pool
	Channel string
}

type meta struct {
	ContentType string   `json:"contentType"`
	HashChain   [][]byte `json:"history"`
	Tags        []string `json:"tags"`
}

// Get returns a library app mounted on the provided path in the pool
func Get(p *pool.Pool, channel string) Library {
	return Library{
		Pool:    p,
		Channel: channel,
	}
}

func (l *Library) updateLocals(documents []Document) []Document {
	for idx := range documents {
		d := &documents[idx]
		if d.Mode&Newer > 0 {
			continue
		}

		stat, err := os.Stat(d.LocalPath)
		if err != nil {
			d.Mode = Deleted
			sqlSetDocument(l.Pool.Name, l.Channel, *d)
		} else if stat.ModTime().Sub(d.LocalModTime) > time.Second {
			d.Mode = Newer
			sqlSetDocument(l.Pool.Name, l.Channel, *d)
		}
	}
	return documents
}

func compareHistory(local, remote *Document) Mode {
	if bytes.Equal(local.Hash, remote.Hash) {
		return Same
	}

	for _, l := range local.HashChain {
		if bytes.Equal(remote.Hash, l) {
			return Older
		}
	}

	for _, r := range remote.HashChain {
		if bytes.Equal(local.Hash, r) {
			return Newer
		}
	}
	return Conflict
}

func (l *Library) updateRemote(documents []Document) []Document {
	m := map[string]*Document{}
	for idx := range documents {
		d := &documents[idx]
		if d.LocalPath != "" {
			m[d.Name] = d
		}
	}

	for idx := range documents {
		d := &documents[idx]
		local := m[d.Name]
		if d.LocalPath != "" || d.Mode&Conflict > 0 || local == nil {
			continue
		}

		if local.Mode&Newer > 0 {
			d.Mode = Conflict
			sqlSetDocument(l.Pool.Name, l.Channel, *d)
			continue
		}

		d.Mode = compareHistory(local, d)
	}
	return documents
}

// List returns the documents in provided folder
func (l *Library) List(folder string) ([]Document, error) {
	hs, _ := l.Pool.List(sqlGetOffset(l.Pool.Name, l.Channel))
	for _, h := range hs {
		l.accept(h)
	}

	folders, err := sqlGetSubfolders(l.Pool.Name, l.Channel, folder)
	if core.IsErr(err, "cannot list subfolders in %s/%s/%s", l.Pool.Name, l.Channel, folder) {
		return nil, err
	}

	documents, err := sqlGetDocumentsInFolder(l.Pool.Name, l.Channel, folder)
	if core.IsErr(err, "cannot list documents in %s/%s/%s", l.Pool.Name, l.Channel, folder) {
		return nil, err
	}

	documents = l.updateLocals(documents)
	documents = l.updateRemote(documents)

	documents = append(documents, folders...)
	return documents, nil
}

func (l *Library) Save(id uint64, dest string) error {
	f, err := os.Create(dest)
	if core.IsErr(err, "cannot create '%s': %v", dest) {
		return err
	}
	defer f.Close()

	err = l.Pool.Receive(id, nil, f)
	if core.IsErr(err, "cannot get file with id %d: %v", id) {
		return err
	}
	return nil
}

func (l *Library) Receive(id uint64, localPath string, tags ...string) (pool.Head, error) {
	f, err := os.Create(localPath)
	if core.IsErr(err, "cannot create '%s': %v", localPath) {
		return pool.Head{}, err
	}
	defer f.Close()

	d, ok, err := sqlGetDocument(l.Pool.Name, l.Channel, id)
	if core.IsErr(err, "cannot get document with id '%d': %v", id) {
		return pool.Head{}, err
	}
	if !ok {
		return pool.Head{}, core.ErrInvalidId
	}

	err = l.Pool.Receive(id, nil, f)
	if core.IsErr(err, "cannot get file with id %d: %v", id) {
		return pool.Head{}, err
	}

	stat, _ := os.Stat(localPath)
	d.LocalModTime = stat.ModTime()
	d.Size = uint64(stat.Size())

	err = sqlSetDocument(l.Pool.Name, l.Channel, d)
	if core.IsErr(err, "cannot update document for id %d: %v", id) {
		return pool.Head{}, err
	}
	return pool.Head{}, nil
}

func (l *Library) Delete(id uint64) error {
	return nil
}

func (l *Library) GetLocalPath(name string) (string, bool) {
	d, ok, _ := sqlGetLocal(l.Pool.Name, l.Channel, name)
	if ok {
		return d.LocalPath, true
	} else {
		return "", false
	}
}

func (l *Library) GetLocalDocument(name string) (Document, bool) {
	d, ok, _ := sqlGetLocal(l.Pool.Name, l.Channel, name)
	return d, ok
}

func (l *Library) Send(localPath string, name string, tags ...string) (pool.Head, error) {
	mime, err := mimetype.DetectFile(localPath)
	if core.IsErr(err, "cannot detect mime type of '%s': %v", localPath) {
		return pool.Head{}, err
	}

	stat, _ := os.Stat(localPath)

	var hashChain [][]byte
	d, ok, err := sqlGetLocal(l.Pool.Name, l.Channel, name)
	if core.IsErr(err, "db error in reading document %s: %v", name) {
		return pool.Head{}, err
	}
	if ok {
		hashChain = append(d.HashChain, d.Hash)
		if len(hashChain) > HashChainMaxLength {
			hashChain = hashChain[len(hashChain)-HashChainMaxLength:]
		}
		if tags == nil {
			tags = d.Tags
		}
	}

	m, err := json.Marshal(meta{
		ContentType: mime.String(),
		Tags:        tags,
		HashChain:   hashChain,
	})
	if core.IsErr(err, "cannot marshal metadata to json: %v") {
		return pool.Head{}, err
	}

	f, err := os.Open(localPath)
	if core.IsErr(err, "cannot open '%s': %v", localPath) {
		return pool.Head{}, err
	}
	defer f.Close()

	h, err := l.Pool.Send(path.Join(l.Channel, name), f, m)
	if core.IsErr(err, "cannot post content to pool '%s': %v", l.Pool.Name) {
		return pool.Head{}, err
	}

	l.Pool.Sync()
	d = Document{
		Id:           h.Id,
		Name:         name,
		Mode:         Same,
		ModTime:      h.ModTime,
		LocalModTime: stat.ModTime(),
		Size:         uint64(h.Size),
		AuthorId:     h.AuthorId,
		ContentType:  mime.String(),
		Hash:         h.Hash,
		HashChain:    hashChain,
		Tags:         tags,
		LocalPath:    localPath,
		Offset:       0,
	}
	err = sqlSetDocument(l.Pool.Name, l.Channel, d)
	if core.IsErr(err, "cannot set document in db for '%s': %v", localPath) {
		return pool.Head{}, err
	}
	return h, err
}

func (l *Library) accept(head pool.Head) {
	if !strings.HasPrefix(head.Name, l.Channel+"/") {
		return
	}

	var m meta
	err := json.Unmarshal(head.Meta, &m)
	if core.IsErr(err, "invalid meta in head: %v") {
		return
	}
	name := head.Name[len(l.Channel)+1:]

	d := Document{
		Id:          head.Id,
		Name:        name,
		ModTime:     head.ModTime,
		Size:        uint64(head.Size),
		AuthorId:    head.AuthorId,
		ContentType: m.ContentType,
		Offset:      head.Offset,
		HashChain:   m.HashChain,
	}

	err = sqlSetDocument(l.Pool.Name, l.Channel, d)
	core.IsErr(err, "cannot save document to db: %v")
}
