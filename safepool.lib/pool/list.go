package pool

import (
	"bytes"
	"encoding/json"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/security"
)

// func (p *Pool) list(prefix string, offset int) ([]Head, error) {
// 	hs, err := sqlGetHeads(p.Name, prefix, offset)
// 	if core.IsErr(err, "cannot read Pool heads: %v") {
// 		return nil, err
// 	}
// 	return hs, err
// }

func (p *Pool) readHead(name string) (Head, error) {
	var b bytes.Buffer
	_, err := p.readFile(name, nil, &b)
	if core.IsErr(err, "cannot read header of %s in %s: %v", name, p.e) {
		return Head{}, err
	}

	var h Head
	err = json.Unmarshal(b.Bytes(), &h)
	if core.IsErr(err, "corrupted header for file %s", name) {
		return Head{}, err
	}

	if !security.Verify(h.AuthorId, h.Hash, h.Signature) {
		return Head{}, ErrNoExchange
	}

	return h, err
}
