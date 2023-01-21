package pool

import (
	"encoding/json"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/code-to-go/safepool.lib/sql"
)

func sqlGetHeads(pool string, offset int) ([]Head, error) {
	rows, err := sql.Query("GET_HEADS", sql.Args{"pool": pool, "offset": offset})
	if core.IsErr(err, "cannot get pools heads from db: %v") {
		return nil, err
	}
	defer rows.Close()

	var heads []Head
	for rows.Next() {
		var h Head
		var modTime int64
		var hash string
		var meta string
		err = rows.Scan(&h.Id, &h.Name, &modTime, &h.Size, &h.AuthorId, &hash, &h.Offset, &meta)
		if !core.IsErr(err, "cannot read pool heads from db: %v") {
			h.Hash = sql.DecodeBase64(hash)
			h.ModTime = sql.DecodeTime(modTime)
			h.Meta = sql.DecodeBase64(meta)
			heads = append(heads, h)
		}
	}
	return heads, nil
}

func sqlGetHead(pool string, id uint64) (Head, error) {
	var h Head
	var modTime int64
	var hash string
	var meta string
	err := sql.QueryRow("GET_HEAD", sql.Args{"pool": pool, "id": id},
		&h.Id, &h.Name, &modTime, &h.Size, &h.AuthorId, &hash, &h.Offset, &meta)
	if core.IsErr(err, "cannot get head with id '%d' in pool '%s': %v", id, pool) {
		return Head{}, err
	}

	h.Hash = sql.DecodeBase64(hash)
	h.ModTime = sql.DecodeTime(modTime)
	h.Meta = sql.DecodeBase64(meta)
	return h, nil
}

func sqlDelHeadBefore(pool string, id uint64) error {
	_, err := sql.Exec("DEL_HEAD_BEFORE", sql.Args{"pool": pool, "beforeId": id})
	return err
}

func sqlAddHead(pool string, h Head) error {
	_, err := sql.Exec("SET_HEAD", sql.Args{
		"pool":     pool,
		"id":       h.Id,
		"name":     h.Name,
		"size":     h.Size,
		"authorId": h.AuthorId,
		"modTime":  sql.EncodeTime(h.ModTime),
		"hash":     sql.EncodeBase64(h.Hash[:]),
		"meta":     sql.EncodeBase64(h.Meta),
	})
	return err
}

func (p *Pool) sqlGetKey(keyId uint64) []byte {
	rows, err := sql.Query("GET_KEY", sql.Args{"pool": p.Name, "keyId": keyId})
	if err != nil {
		return nil
	}
	defer rows.Close()

	for rows.Next() {
		var key string
		var err = rows.Scan(&key)
		if !core.IsErr(err, "cannot read key from db: %v") {
			return sql.DecodeBase64(key)
		}
	}
	return nil
}

func (p *Pool) sqlSetKey(keyId uint64, value []byte) error {
	_, err := sql.Exec("SET_KEY", sql.Args{"pool": p.Name, "keyId": keyId, "keyValue": sql.EncodeBase64(value)})
	return err
}

func (p *Pool) sqlGetKeystore() (Keystore, error) {
	rows, err := sql.Query("GET_KEYS", sql.Args{"pool": p.Name})
	if core.IsErr(err, "cannot read keystore for pool %s: %v", p.Name) {
		return nil, err
	}
	defer rows.Close()

	ks := Keystore{}
	for rows.Next() {
		var keyId uint64
		var keyValue string
		var err = rows.Scan(&keyId, &keyValue)
		if !core.IsErr(err, "cannot read key from db: %v") {
			ks[keyId] = sql.DecodeBase64(keyValue)
		}
	}
	return ks, nil
}

func (p *Pool) sqlGetAccesses(onlyTrusted bool) (identities []security.Identity, accesses []Access, err error) {
	var q string
	if onlyTrusted {
		q = "GET_TRUSTED_ACCESSES"
	} else {
		q = "GET_ACCESSES"
	}

	rows, err := sql.Query(q, sql.Args{"pool": p.Name})
	if core.IsErr(err, "cannot get trusted identities from db: %v") {
		return nil, nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var i security.Identity
		var id string
		var i64 string
		var state State
		var modTime int64
		var ts int64
		err = rows.Scan(&id, &i64, &state, &modTime, &ts)
		if core.IsErr(err, "cannot read identity from db: %v") {
			continue
		}

		i, err = security.IdentityFromBase64(i64)
		if core.IsErr(err, "invalid identity '%s': %v", i64) {
			continue
		}
		identities = append(identities, i)

		accesses = append(accesses, Access{
			Id:      id,
			ModTime: sql.DecodeTime(modTime),
			State:   state,
		})
	}
	return identities, accesses, nil
}

func (p *Pool) sqlSetAccess(a Access) error {
	_, err := sql.Exec("SET_ACCESS", sql.Args{
		"id":      a.Id,
		"pool":    p.Name,
		"modTime": sql.EncodeTime(a.ModTime),
		"state":   a.State,
		"ts":      sql.EncodeTime(core.Now()),
	})
	return err
}

func sqlSave(name string, c Config) error {
	data, err := json.Marshal(&c)
	if core.IsErr(err, "cannot marshal transport configuration of %s: %v", name) {
		return err
	}

	_, err = sql.Exec("SET_POOL", sql.Args{"name": name, "configs": sql.EncodeBase64(data)})
	core.IsErr(err, "cannot save transport configuration of %s: %v", name)
	return err
}

func sqlLoad(name string) (Config, error) {
	var blob string
	var c Config
	err := sql.QueryRow("GET_POOL", sql.Args{"name": name}, &blob)
	if core.IsErr(err, "cannot get pool %s config: %v", name) {
		return Config{}, err
	}

	data := sql.DecodeBase64(blob)
	err = json.Unmarshal(data, &c)
	core.IsErr(err, "cannot unmarshal configs of %s: %v", name)
	return c, err
}

func sqlList() ([]string, error) {
	var names []string
	rows, err := sql.Query("LIST_POOL", nil)
	if core.IsErr(err, "cannot list pools: %v") {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var n string
		err = rows.Scan(&n)
		if err == nil {
			names = append(names, n)
		}
	}
	return names, err
}
