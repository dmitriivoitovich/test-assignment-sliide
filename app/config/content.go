package config

import "github.com/dmitriivoitovich/test-assignment-sliide/app/provider"

type ContentMix []ContentConfig

type ContentConfig struct {
	Type     provider.Provider
	Fallback *provider.Provider
}

var (
	config1 = ContentConfig{
		Type:     provider.Provider1,
		Fallback: &provider.Provider2,
	}
	config2 = ContentConfig{
		Type:     provider.Provider2,
		Fallback: &provider.Provider3,
	}
	config3 = ContentConfig{
		Type:     provider.Provider3,
		Fallback: &provider.Provider1,
	}
	config4 = ContentConfig{
		Type:     provider.Provider1,
		Fallback: nil,
	}

	DefaultContentMix = ContentMix{config1, config1, config2, config3, config4, config1, config1, config2}
)
