package storage

import (
	"github.com/dollarkillerx/ocean/pkg/models"

	"sync"
)

type Storage struct {
	Schema       map[string]models.Schema            // schema
	ListData     map[string][]map[string]interface{} // 具体数据
	GlobalLock   sync.RWMutex                        // 全局锁 用于raft 同步
	SchemaRWLock map[string]sync.RWMutex             // schema 局部锁
	DataRWLock   map[string]sync.RWMutex             // 数据局部锁
}

func New() *Storage {
	return &Storage{
		Schema:       map[string]models.Schema{},
		ListData:     map[string][]map[string]interface{}{},
		SchemaRWLock: map[string]sync.RWMutex{},
		DataRWLock:   map[string]sync.RWMutex{},
	}
}
