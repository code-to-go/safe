package pool

import (
	"sync"

	"github.com/code-to-go/safe/safepool/transport"
)

var Connections = map[string]transport.Exchanger{}
var ConnectionsMutex = &sync.Mutex{}

const pingName = ".reserved.ping.%d.test"
