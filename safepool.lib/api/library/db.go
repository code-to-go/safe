package library

import (
	"encoding/json"
	"path"
	"strings"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/code-to-go/safepool.lib/sql"
)

func sqlSetDocument(pool string, base string, d Document) error {
	folder, level := getFolderAndLevel(path.Dir(d.Name))
	hashChain, _ := json.Marshal(d.HashChain)

	_, err := sql.Exec("SET_DOCUMENT", sql.Args{"pool": pool, "base": base, "id": d.Id, "name": d.Name,
		"authorId": d.AuthorId, "modTime": sql.EncodeTime(d.ModTime), "mode": d.Mode, "size": d.Size,
		"contentType": d.ContentType, "hash": sql.EncodeBase64(d.Hash), "hashChain": hashChain,
		"localModTime": sql.EncodeTime(d.LocalModTime), "localPath": d.LocalPath,
		"offset": d.Offset, "folder": folder, "level": level})
	if core.IsErr(err, "cannot set document %d on db: %v", d.Name) {
		return err
	}

	if d.LocalPath != "" {
		_, err = sql.Exec("CLEAN_DOCUMENT_LOCAL", sql.Args{"pool": pool, "base": base, "id": d.Id, "name": d.Name})
		if core.IsErr(err, "cannot clean document local for '%d' on db: %v", d.Name) {
			return err
		}
	}

	return err
}

func sqlGetDocument(pool string, base string, id uint64) (Document, bool, error) {
	var d Document
	var attrs []byte

	err := sql.QueryRow("GET_DOCUMENT", sql.Args{"pool": pool, "base": base, "id": id}, &attrs)
	if err == sql.ErrNoRows {
		return Document{}, false, nil
	} else if core.IsErr(err, "cannot get document '%d' from DB: %v", id) {
		return d, false, err
	}

	err = json.Unmarshal(attrs, &d)
	if core.IsErr(err, "cannot unmarshal document '%d' from DB: %v", id) {
		return d, false, err
	}

	return d, true, err
}

func sqlGetLocal(pool, base, name string) (Document, bool, error) {
	var d Document
	var hash string
	var hashChain []byte
	var modTime int64
	var localModTime int64

	err := sql.QueryRow("GET_DOCUMENT_LOCAL", sql.Args{"pool": pool, "base": base, "name": name},
		&d.Name, &d.AuthorId, &d.Mode, &modTime, &d.Id, &d.Size, &d.ContentType, &hash, &hashChain,
		&localModTime, &d.LocalPath, &d.Offset)
	if err == sql.ErrNoRows {
		return Document{}, false, nil
	} else if core.IsErr(err, "cannot get local document '%d' from DB: %v", name) {
		return d, false, err
	}

	d.Hash = sql.DecodeBase64(hash)
	json.Unmarshal(hashChain, &d.HashChain)
	d.ModTime = sql.DecodeTime(modTime)
	d.LocalModTime = sql.DecodeTime(localModTime)

	return d, true, nil
}

// func sqlSetlocal(pool, base, name string, id uint64, localPath string, modTime time.Time, size ) error {
// 	_, err := sql.Exec("SET_DOCUMENT_LEAD", sql.Args{"pool": pool, "base": base, "name": name,
// 		"localPath": localPath, "id": id})
// 	return err
// }

func getFolderAndLevel(folder string) (string, int) {
	folder = path.Clean(folder)
	folder = strings.TrimPrefix(folder, ".")
	folder = strings.Trim(folder, "/")
	level := strings.Count(folder, "/")
	return folder, level
}

func sqlGetSubfolders(pool string, base string, folder string) ([]Document, error) {
	folder, level := getFolderAndLevel(folder)
	rows, err := sql.Query("GET_DOCUMENTS_SUBFOLDERS", sql.Args{"pool": pool, "base": base, "folder": folder + "%", "level": level + 1})
	if core.IsErr(err, "cannot query documents from db: %v") {
		return nil, err
	}
	var documents []Document
	for rows.Next() {
		d := Document{
			Mode: Folder,
		}
		err = rows.Scan(&d.Name)
		if !core.IsErr(err, "cannot scan row in Documents: %v", err) {
			documents = append(documents, d)
		}
	}
	return documents, nil
}

//name,authorId,mode,id,size,contentType,hash,hashChain,localPath,offset

func sqlGetDocumentsInFolder(pool string, base string, folder string) ([]Document, error) {
	rows, err := sql.Query("GET_DOCUMENTS_IN_FOLDER", sql.Args{"pool": pool, "base": base, "folder": folder})
	if core.IsErr(err, "cannot query documents from db: %v") {
		return nil, err
	}
	var documents []Document
	for rows.Next() {
		var d Document
		var hash string
		var hashChain []byte
		var modTime int64
		var localModTime int64

		err = rows.Scan(&d.Name, &d.AuthorId, &d.Mode, &modTime, &d.Id, &d.Size, &d.ContentType, &hash, &hashChain,
			&localModTime, &d.LocalPath, &d.Offset)
		if !core.IsErr(err, "cannot scan row in Documents: %v", err) {
			d.Hash = sql.DecodeBase64(hash)
			json.Unmarshal(hashChain, &d.HashChain)
			d.ModTime = sql.DecodeTime(modTime)
			d.LocalModTime = sql.DecodeTime(localModTime)
			documents = append(documents, d)
		}
	}
	return documents, nil
}

func sqlGetOffset(pool string, base string) int {
	var offset int
	err := sql.QueryRow("GET_DOCUMENTS_OFFSET", sql.Args{"pool": pool, "base": base}, &offset)
	if err == nil {
		return offset
	} else {
		return -1
	}
}
