package pool

import (
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/godruoyi/go-snowflake"
)

const FeedsFolder = "feeds"

func (p *Pool) getSlots() ([]string, error) {
	slot, _ := sqlGetSlot(p.Name, p.e.String())

	fs, err := p.e.ReadDir(path.Join(p.Name, FeedsFolder), 0)
	if core.IsErr(err, "cannot list slots in '%v': %v", p) {
		return nil, err
	}

	var slots []string
	for _, f := range fs {
		if f.Name() >= slot {
			slots = append(slots, slot)
		}
	}

	sort.Strings(slots)
	return slots, nil
}

func (p *Pool) Sync() error {
	if !p.e.Touched(p.Name + "/") {
		return nil
	}
	hs, _ := p.List(0)

	feeds := map[uint64]Feed{}
	for _, h := range hs {
		feeds[h.Id] = h
	}

	slots, err := p.getSlots()
	if err != nil {
		return err
	}

	thresold := uint64(core.Since(core.SnowFlakeStart) - LifeSpan)
	for _, slot := range slots {
		fs, err := p.e.ReadDir(path.Join(p.Name, FeedsFolder), 0)
		if core.IsErr(err, "cannot read content in slot %s in pool", slot, p) {
			continue
		}
		for _, f := range fs {
			name := f.Name()
			if !strings.HasSuffix(name, ".feed") {
				continue
			}

			id, err := strconv.ParseInt(name[0:len(name)-len(".feed")], 10, 64)
			if err != nil {
				continue
			}
			if _, found := feeds[uint64(id)]; found {
				continue
			}

			sid := snowflake.ParseID(uint64(id))
			if sid.Timestamp < thresold {
				continue
			}

			f, err := p.readHead(path.Join(p.Name, name))
			if core.IsErr(err, "cannot read file %s from %s: %v", name, p.e) {
				continue
			}
			f.Slot = slot
			_ = sqlAddFeed(p.Name, f)
			hs = append(hs, f)
		}
		sqlSetSlot(p.Name, p.e.String(), slot)
	}

	if time.Until(p.lastReplica) > ReplicaPeriod {
		p.replica()
		p.lastReplica = core.Now()
	}
	return nil
}
