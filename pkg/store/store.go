package store

import (
	"sort"
	"sync"

	"github.com/leyou240/speedle-plus/api/pms"
	"github.com/leyou240/speedle-plus/pkg/errors"
)

var (
	storeBuildersMu = &sync.RWMutex{}
	storeBuilders   = make(map[string]StoreBuilder)
)

type StoreBuilder interface {
	NewStore(storeConfig map[string]interface{}) (pms.PolicyStoreManager, error)
	GetStoreParams() map[string]string
}

// Register makes a type of store available by the provided name.
// If Register is called twice with the same name or if storeBuilder is nil,
// it panics.
func Register(storeType string, storeBuilder StoreBuilder) {
	storeBuildersMu.Lock()
	defer storeBuildersMu.Unlock()
	if storeBuilder == nil {
		panic("speedle: Register storeBuilder is nil")
	}
	if _, dup := storeBuilders[storeType]; dup {
		panic("speedle: Register called twice for storeBuilder " + storeType)
	}
	storeBuilders[storeType] = storeBuilder
}

func unregisterAllStoreBuilders() {
	storeBuildersMu.Lock()
	defer storeBuildersMu.Unlock()
	// For tests.
	storeBuilders = make(map[string]StoreBuilder)
}

// StoreBuilders returns a sorted list of the names of the registered store types.
func StoreBuilders() []string {
	storeBuildersMu.RLock()
	defer storeBuildersMu.RUnlock()
	var list []string
	for storeType := range storeBuilders {
		list = append(list, storeType)
	}
	sort.Strings(list)
	return list
}

// GetAllStoreParams The return map expect param's command line flag name as key, and storeProps key in config file as value
func GetAllStoreParams() map[string]string {
	storeBuildersMu.RLock()
	defer storeBuildersMu.RUnlock()
	paramMap := make(map[string]string)

	for _, storeBuilder := range storeBuilders {
		for k, v := range storeBuilder.GetStoreParams() {
			paramMap[k] = v
		}
	}
	return paramMap
}

func NewStore(storeType string, storeConfig map[string]interface{}) (pms.PolicyStoreManager, error) {
	storeBuildersMu.RLock()
	storeBuilder, ok := storeBuilders[storeType]
	storeBuildersMu.RUnlock()
	if !ok {
		return nil, errors.Errorf(errors.ConfigError, "unknown store type %q (forgotten import?)", storeType)
	}
	return storeBuilder.NewStore(storeConfig)
}
