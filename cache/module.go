package cache

import (
	"github.com/Ashishkapoor1469/Nestgo/common"
	"github.com/Ashishkapoor1469/Nestgo/di"
	"github.com/Ashishkapoor1469/Nestgo/middleware"
)

// CacheModule registers a CacheStore provider into the DI container.
type CacheModule struct {
	store middleware.CacheStore
}

// NewCacheModule creates a new CacheModule.
// If no store is provided, it defaults to MemoryCacheStore.
func NewCacheModule(store ...middleware.CacheStore) *CacheModule {
	var s middleware.CacheStore
	if len(store) > 0 {
		s = store[0]
	} else {
		s = middleware.NewMemoryCacheStore()
	}
	return &CacheModule{store: s}
}

// Module returns the NestGo module configuration.
func (m *CacheModule) Module() common.ModuleConfig {
	return common.ModuleConfig{
		Name: "CacheModule",
		Providers: []di.Provider{
			{
				Instance: m.store,
			},
		},
	}
}
