package nosql

import (
	"strconv"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/syndtr/goleveldb/leveldb"
)

var manager *Manager

func init() {
	_ = GetManager()
}

func GetManager() *Manager {
	if manager == nil {
		manager = &Manager{
			RedisConnections:   make(map[string]*redisClientHolder),
			LevelDBConnections: make(map[string]*levelDBHolder),
		}
	}
	return manager
}

// Manager is the nosql connection manager
type Manager struct {
	mutex sync.Mutex

	RedisConnections   map[string]*redisClientHolder
	LevelDBConnections map[string]*levelDBHolder
}

type redisClientHolder struct {
	redis.UniversalClient
	name  []string
	count int64
}

type levelDBHolder struct {
	name  []string
	count int64
	db    *leveldb.DB
}

func (r *redisClientHolder) Close() error {
	return manager.CloseRedisClient(r.name[0])
}

func valToTimeDuration(vs []string) (result time.Duration) {
	var err error
	for _, v := range vs {
		result, err = time.ParseDuration(v)
		if err != nil {
			var val int
			val, err = strconv.Atoi(v)
			result = time.Duration(val)
		}
		if err == nil {
			return
		}
	}
	return
}
