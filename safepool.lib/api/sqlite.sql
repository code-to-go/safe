-- INIT
CREATE TABLE IF NOT EXISTS identities (
    id VARCHAR(256),
    i64 BLOB,
    trusted INTEGER,
    alias VARCHAR(256),
    PRIMARY KEY(id)
);

-- INIT
CREATE INDEX IF NOT EXISTS idx_identities_trust ON identities(trusted);

-- GET_IDENTITIES
SELECT i64, alias FROM identities

-- GET_IDENTITY
SELECT i64, alias FROM identities WHERE id=:id

-- GET_TRUSTED
SELECT i64, alias FROM identities WHERE trusted

-- SET_TRUSTED
UPDATE identities SET trusted=:trusted WHERE id=:id

-- SET_IDENTITY
INSERT INTO identities(id,i64,alias) VALUES(:id,:i64,'')
    ON CONFLICT(id) DO UPDATE SET i64=:i64
	WHERE id=:id

-- SET_ALIAS
UPDATE identities SET alias=:alias WHERE id=:id

-- INIT
CREATE TABLE IF NOT EXISTS configs (
    pool VARCHAR(128) NOT NULL, 
    k VARCHAR(64) NOT NULL, 
    s VARCHAR(64) NOT NULL,
    i INTEGER NOT NULL,
    b TEXT,
    CONSTRAINT pk_safe_key PRIMARY KEY(pool,k)
);

-- GET_CONFIG
SELECT s, i, b FROM configs WHERE pool=:pool AND k=:key

-- SET_CONFIG
INSERT INTO configs(pool,k,s,i,b) VALUES(:pool,:key,:s,:i,:b)
	ON CONFLICT(pool,k) DO UPDATE SET s=:s,i=:i,b=:b
	WHERE pool=:pool AND k=:key

-- INIT
CREATE TABLE IF NOT EXISTS heads (
    offset INTEGER PRIMARY KEY AUTOINCREMENT,
    pool VARCHAR(128) NOT NULL, 
    id INTEGER NOT NULL,
    name VARCHAR(8192) NOT NULL, 
    modTime INTEGER NOT NULL,
    size INTEGER NOT NULL,
    authorId VARCHAR(80) NOT NULL,
    hash VARCHAR(128) NOT NULL, 
    meta VARCHAR(4096) NOT NULL
)

-- INIT
CREATE INDEX IF NOT EXISTS idx_heads_id ON heads(id);

-- INIT
CREATE INDEX IF NOT EXISTS idx_heads_pool ON heads(pool);

-- INIT
CREATE INDEX IF NOT EXISTS idx_heads_name ON heads(name);

-- GET_HEADS
SELECT id, name, modTime, size, authorId, hash, offset, meta FROM heads WHERE pool=:pool AND offset > :offset ORDER BY offset

-- GET_HEAD
SELECT id, name, modTime, size, authorId, hash, offset, meta FROM heads WHERE pool=:pool AND id=:id

-- SET_HEAD
INSERT INTO heads(pool,id,name,modTime,size,authorId,hash,meta) VALUES(:pool,:id,:name,:modTime,:size,:authorId,:hash,:meta)

-- INIT
CREATE TABLE IF NOT EXISTS keys (
    pool VARCHAR(128) NOT NULL, 
    keyId INTEGER, 
    keyValue VARCHAR(128),
    CONSTRAINT pk_safe_keyId PRIMARY KEY(pool,keyId)
);

-- GET_KEYS
SELECT keyId, keyValue FROM keys WHERE pool=:pool

-- GET_KEY
SELECT keyValue FROM keys WHERE pool=:pool AND keyId=:keyId

-- SET_KEY
INSERT INTO keys(pool,keyId,keyValue) VALUES(:pool,:keyId,:keyValue)
    ON CONFLICT(pool,keyId) DO UPDATE SET keyValue=:keyValue
	    WHERE pool=:pool AND keyId=:keyId

-- INIT
CREATE TABLE IF NOT EXISTS pools (
    name VARCHAR(128),
    configs BLOB,
    PRIMARY KEY(name)
);

-- GET_POOL
SELECT configs FROM pools WHERE name=:name

-- LIST_POOL
SELECT name FROM pools

-- SET_POOL
INSERT INTO pools(name,configs) VALUES(:name,:configs)
    ON CONFLICT(name) DO UPDATE SET configs=:configs
	    WHERE name=:name

-- INIT
CREATE TABLE IF NOT EXISTS accesses (
    pool VARCHAR(128),
    id VARCHAR(256),
    state INTEGER,
    modTime INTEGER,
    ts INTEGER,
    CONSTRAINT pk_safe_sig_enc PRIMARY KEY(pool,id)
);

-- GET_TRUSTED_ACCESSES
SELECT s.id, i.i64, state, modTime, ts FROM identities i INNER JOIN accesses s WHERE s.pool=:pool AND (i.id = s.id OR i.id IS NULL) AND i.trusted

-- GET_ACCESSES
SELECT s.id, i.i64, state, modTime, ts FROM identities i INNER JOIN accesses s WHERE s.pool=:pool AND (i.id = s.id OR i.id IS NULL)

-- GET_ACCESS
SELECT state, modTime, ts FROM accesses s WHERE s.pool=:pool AND id = :id 

-- SET_ACCESS
INSERT INTO accesses(pool,id,state,modTime,ts) VALUES(:pool,:id,:state,:modTime,:ts)
    ON CONFLICT(pool,id) DO UPDATE SET state=:state,modTime=:modTime,ts=:ts WHERE
    pool=:pool AND id=:id

-- DEL_GRANT
DELETE FROM accesses WHERE id=:id AND pool=:pool

-- INIT
CREATE TABLE IF NOT EXISTS chats (
    pool VARCHAR(128),
    id INTEGER,
    author VARCHAR(128),
    message BLOB,
    offset INTEGER,
    CONSTRAINT pk_pool_id_author PRIMARY KEY(pool,id,author)
);

-- SET_CHAT_MESSAGE
INSERT INTO chats(pool,id,author,message,offset) VALUES(:pool,:id,:author,:message, :offset)
    ON CONFLICT(pool,id,author) DO UPDATE SET message=:message
	    WHERE pool=:pool AND id=:id AND author=:author

-- GET_CHAT_MESSAGES
SELECT message FROM chats WHERE pool=:pool AND id > :afterId AND id < :beforeId ORDER BY id DESC LIMIT :limit

-- GET_CHATS_OFFSET
SELECT max(offset) FROM chats WHERE pool=:pool

-- INIT
CREATE TABLE IF NOT EXISTS documents (
    pool VARCHAR(128) NOT NULL,
    base VARCHAR(128) NOT NULL,
    id INTEGER NOT NULL,
    name VARCHAR(4096) NOT NULL,
    authorId VARCHAR(128) NOT NULL,
    mode INTEGER NOT NULL,
    modTime INTEGER,
    size INTEGER NOT NULL,
    contentType VARCHAR(128) NOT NULL,
    hash VARCHAR(128) NOT NULL,
    hashChain BLOB,
    localModTime INTEGER,
    localPath VARCHAR(4096) NOT NULL,
    offset INTEGER NOT NULL,
    folder VARCHAR(4096) NOT NULL,
    level INTEGER NOT NULL,
    CONSTRAINT pk_pool_base_id PRIMARY KEY(pool,base,name,authorId,localPath)
);

-- INIT
CREATE INDEX IF NOT EXISTS idx_documents_name ON heads(name);

-- SET_DOCUMENT
INSERT INTO documents(pool,base,id,name,authorId,mode,modTime,size,contentType,hash,hashChain,localModTime,localPath,offset,folder,level) 
    VALUES(:pool,:base,:id,:name,:authorId,:mode,:modTime,:size,:contentType,:hash,:hashChain,:localModTime,:localPath,:offset,:folder,:level)
    ON CONFLICT(pool,base,name,authorId,localPath) DO UPDATE SET id=:id,mode=:mode,modTime=:modTime,size=:size,
    contentType=:contentType,hash=:hash, hashChain=:hashChain,localModTime=:localModTime,localPath=:localPath,offset=:offset
	    WHERE pool=:pool AND base=:base AND name=:name AND authorId=:authorId AND localPath=:localPath

-- CLEAN_DOCUMENT_LOCAL
UPDATE documents SET localPath='' WHERE pool=:pool AND base=:base AND name=:name AND id!=:id;

-- SET_DOCUMENT_MODE
UPDATE documents SET mode=:mode WHERE pool=:pool AND base=:base AND id=:id;

-- GET_DOCUMENT
SELECT name,authorId,mode,modTime,id,size,contentType,hash,hashChain,localModTime,localPath,offset FROM documents 
    WHERE pool=:pool AND base=:base AND id=:id

-- GET_DOCUMENT_LOCAL
SELECT name,authorId,mode,modTime,id,size,contentType,hash,hashChain,localModTime,localPath,offset FROM documents 
    WHERE pool=:pool AND base=:base AND name=:name AND localPath != ""

-- GET_DOCUMENTS_IN_FOLDER
SELECT name,authorId,mode,modTime,id,size,contentType,hash,hashChain,localModTime,localPath,offset FROM documents 
    WHERE pool=:pool AND base=:base AND folder=:folder ORDER BY name

-- GET_DOCUMENTS_SUBFOLDERS
SELECT folder FROM documents WHERE pool=:pool AND base=:base AND folder LIKE :folder AND level=:level ORDER BY folder

-- GET_DOCUMENTS_OFFSET
SELECT max(offset) FROM documents WHERE pool=:pool AND base=:base 
