package providers

import (
	"context"
	"sync"
)

type Provider interface {
	Start(ctx context.Context) ([]string, error)
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

func Start(ctx context.Context, provider string) ([]string, error) {
	if val, ok := Providers[provider]; ok {
		links, err := val.Provider.Start(ctx)

		if err != nil {
			return nil, err
		}

		return links, nil
	}

	return nil, nil
}
