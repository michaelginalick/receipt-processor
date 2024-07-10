package inmemorydb

import "sync"

type Client interface {
	Save(key string, value any)
	Get(key string) (value any, ok bool)
}

type inMemoryDB struct {
	syncMap sync.Map
}

func NewClient() Client {
	return &inMemoryDB{syncMap: sync.Map{}}
}

func (inMem *inMemoryDB) Save(key string, value any) {
	inMem.syncMap.Store(key, value)
}

func (inMem *inMemoryDB) Get(key string) (value any, ok bool) {
	return inMem.syncMap.Load(key)
}
