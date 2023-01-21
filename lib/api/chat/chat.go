package chat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	pool "github.com/code-to-go/safepool.lib/pool"
	"github.com/code-to-go/safepool.lib/security"
	"github.com/godruoyi/go-snowflake"
	"github.com/sirupsen/logrus"
)

type Message struct {
	Id          uint64    `json:"id,string"`
	Author      string    `json:"author"`
	Time        time.Time `json:"time"`
	Content     string    `json:"content"`
	ContentType string    `json:"contentType"`
	Attachments [][]byte  `json:"attachments"`
	Signature   []byte    `json:"signature"`
}

func getHash(m *Message) []byte {
	h := security.NewHash()
	h.Write([]byte(m.Content))
	h.Write([]byte(m.ContentType))
	h.Write([]byte(m.Author))
	for _, a := range m.Attachments {
		h.Write(a)
	}
	return h.Sum(nil)
}

type Chat struct {
	Pool    *pool.Pool
	Channel string
}

func Get(p *pool.Pool, channel string) Chat {
	return Chat{
		Pool:    p,
		Channel: channel,
	}
}

func (c *Chat) TimeOffset(s *pool.Pool) int {
	return sqlGetOffset(s.Name)
}

func (c *Chat) Accept(s *pool.Pool, head pool.Head) bool {
	name := head.Name
	if !strings.HasPrefix(name, "/chat/") || !strings.HasSuffix(name, ".chat") || head.Size > 10*1024*1024 {
		return false
	}
	name = path.Base(name)
	id, err := strconv.ParseInt(name[0:len(name)-5], 10, 64)
	if err != nil {
		return false
	}

	buf := bytes.Buffer{}
	err = s.Receive(head.Id, nil, &buf)
	if core.IsErr(err, "cannot read %s from %s: %v", head.Name, s.Name) {
		return true
	}

	var m Message
	err = json.Unmarshal(buf.Bytes(), &m)
	if core.IsErr(err, "invalid chat message %s: %v", head.Name) {
		return true
	}

	h := getHash(&m)
	if !security.Verify(m.Author, h, m.Signature) {
		logrus.Error("message %s has invalid signature", head.Name)
		return true
	}

	err = sqlSetMessage(s.Name, uint64(id), m.Author, m, head.Offset)
	core.IsErr(err, "cannot write message %s to db:%v", head.Name)
	return true
}

func (c *Chat) SendMessage(content string, contentType string, attachments [][]byte) (uint64, error) {
	m := Message{
		Id:          snowflake.ID(),
		Author:      c.Pool.Self.Id(),
		Time:        core.Now(),
		Content:     content,
		ContentType: contentType,
		Attachments: attachments,
	}
	h := getHash(&m)
	signature, err := security.Sign(c.Pool.Self, h)
	if core.IsErr(err, "cannot sign chat message: %v") {
		return 0, err
	}
	m.Signature = signature

	data, err := json.Marshal(m)
	if core.IsErr(err, "cannot sign chat message: %v") {
		return 0, err
	}

	go func() {
		name := fmt.Sprintf("%s/%d.chat", c.Channel, m.Id)
		_, err = c.Pool.Send(name, bytes.NewBuffer(data), nil)
		core.IsErr(err, "cannot write chat message: %v")
	}()

	core.Info("added chat message with id %d", m.Id)
	return m.Id, nil
}

func (c *Chat) accept(h pool.Head) {
	name := h.Name
	if !strings.HasPrefix(name, c.Channel) || !strings.HasSuffix(name, ".chat") || h.Size > 10*1024*1024 {
		return
	}
	name = path.Base(name)
	id, err := strconv.ParseInt(name[0:len(name)-5], 10, 64)
	if err != nil {
		return
	}

	buf := bytes.Buffer{}
	err = c.Pool.Receive(h.Id, nil, &buf)
	if core.IsErr(err, "cannot read %s from %s: %v", h.Name, h.Name) {
		return
	}

	var m Message
	err = json.Unmarshal(buf.Bytes(), &m)
	if core.IsErr(err, "invalid chat message %s: %v", h.Name) {
		return
	}

	hash := getHash(&m)
	if !security.Verify(m.Author, hash, m.Signature) {
		logrus.Error("message %s has invalid signature", h.Name)
		return
	}

	err = sqlSetMessage(c.Pool.Name, uint64(id), m.Author, m, h.Offset)
	core.IsErr(err, "cannot write message %s to db:%v", h.Name)
}

func (c *Chat) GetMessages(afterId, beforeId uint64, limit int) ([]Message, error) {
	c.Pool.Sync()
	hs, err := c.Pool.List(sqlGetOffset(c.Pool.Name))
	for _, h := range hs {
		c.accept(h)
	}

	messages, err := sqlGetMessages(c.Pool.Name, afterId, beforeId, limit)
	if core.IsErr(err, "cannot read messages from db: %v") {
		return nil, err
	}

	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Id < messages[j].Id
	})
	return messages, nil
}