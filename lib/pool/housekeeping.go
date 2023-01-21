package pool

import (
	"math/rand"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/code-to-go/safepool.lib/core"
	"github.com/godruoyi/go-snowflake"
)

// LifeSpan is the maximal time data should stay in the pool. It is default to 30 days and it cannot be less than 7.
var LifeSpan = 30 * 24 * time.Hour

const sevenDays = 7 * 24 * time.Hour

func (p *Pool) startHouseKeeping() {
	if LifeSpan < sevenDays {
		LifeSpan = sevenDays
	}

	LifeSpan = time.Hour
	p.houseKeeping = time.NewTicker(time.Hour)
	p.stopHouseKeeping = make(chan bool)

	go func() {
		rand.Seed(time.Now().UnixNano())
		n := rand.Intn(600)
		time.Sleep(time.Duration(n) * time.Second)

		for {
			p.HouseKeeping()
			select {
			case <-p.stopHouseKeeping:
				return
			case <-p.houseKeeping.C:
				continue
			}
		}
	}()
}

// HouseKeeping removes old files from the pool. It is automatically called once a day and there is not need to call programmatically
func (p *Pool) HouseKeeping() {
	p.houseKeepingLock.Lock()
	defer p.houseKeepingLock.Unlock()

	thresold := uint64(core.Since(core.SnowFlakeStart) - LifeSpan)
	for _, e := range p.exchangers {
		fs, err := e.ReadDir(path.Join(p.Name, FeedsFolder), 0)
		if core.IsErr(err, "cannot read content in pool %s/%s", e, p.Name) {
			continue
		}
		for _, f := range fs {
			name := f.Name()
			if !strings.HasSuffix(name, ".head") {
				continue
			}

			id_, err := strconv.ParseInt(name[0:len(name)-len(".head")], 10, 64)
			if err != nil {
				continue
			}

			id := uint64(id_)
			sid := snowflake.ParseID(id)
			if sid.Timestamp < thresold {
				p.e.Delete(path.Join(p.Name, FeedsFolder, name))
			}
		}
	}
	thresoldId := thresold << (snowflake.SequenceLength + snowflake.MachineIDLength)
	sqlDelHeadBefore(p.Name, thresoldId)
}
