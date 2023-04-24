package providers

import (
	"Scythe/core/common"
	"context"
	"sync"
)

type Provider interface {
	Start(ctx context.Context) ([]common.Media, error)
}

type ProviderInfo struct {
	Category string
	Provider Provider
}

var mutable sync.Mutex
var Providers = make(map[string]ProviderInfo)

func Register(name, category string, p Provider) {
	mutable.Lock()
	Providers[name] = ProviderInfo{
		Category: category,
		Provider: p,
	}
	mutable.Unlock()
}

func ListProvidersByCategory(category string) []string {
	var providers []string
	for name, val := range Providers {
		if val.Category == category {
			providers = append(providers, name)
		}
	}

	return providers
}

func Start(ctx context.Context, provider string) ([]common.Media, error) {
	if val, ok := Providers[provider]; ok {
		media, err := val.Provider.Start(ctx)

		if err != nil {
			return nil, err
		}

		return media, nil
	}

	return nil, nil
}
