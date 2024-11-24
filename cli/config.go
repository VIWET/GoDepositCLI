package cli

// IndexedConfig stores values by key index
type IndexedConfig[T any] struct {
	Config map[uint32]T `json:"config"`
}

// Get value by key index
func (cfg *IndexedConfig[T]) Get(index uint32) (T, bool) {
	value, ok := cfg.Config[index]
	return value, ok
}

// IndexedConfigWithDefault stroes values by key index and default value
type IndexedConfigWithDefault[T any] struct {
	Default T `json:"default"`
	IndexedConfig[T]
}

// Get value by key index or default
func (cfg *IndexedConfigWithDefault[T]) Get(index uint32) T {
	value, ok := cfg.IndexedConfig.Get()
	if ok {
		return value
	}
	return cfg.Default
}
