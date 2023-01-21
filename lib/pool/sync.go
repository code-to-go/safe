package pool

import (
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/godruoyi/go-snowflake"
)

const FeedsFolder = "feeds"

func (p *Pool) Sync() error {
	if !p.e.Touched(p.Name + "/") {
		return nil
	}
	hs, _ := p.List(0)

	heads := map[uint64]Head{}
	for _, h := range hs {
		heads[h.Id] = h
	}

	fs, err := p.e.ReadDir(path.Join(p.Name, FeedsFolder), 0)
	if core.IsErr(err, "cannot read content in pool %s/%s", p.e, p.Name) {
		return err
	}
	thresold := uint64(core.Since(core.SnowFlakeStart) - LifeSpan)
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

		sid := snowflake.ParseID(uint64(id))
		if sid.Timestamp < thresold {
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
		p.lastReplica = core.Now()
	}
	return nil
}
