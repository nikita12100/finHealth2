package common

import (
	"log/slog"
	"time"
)

type Cache[T any] struct {
	Value   T         `json:"value"`
	Created time.Time `json:"created"`
}

func Cached[K any, V any](
	key K,
	ttl time.Duration,
	dbGet func(key K) (Cache[V], error),
	getValue func(key K) (V, error),
	dbSave func(key K, newValue Cache[V]) error,
) (V, error) {
	var cacheValue V
	if entry, err := dbGet(key); err == nil {
		cacheValue = entry.Value
		if time.Since(entry.Created) < ttl {
			slog.Debug("Used cache", "key", key, "cacheValue", cacheValue)
			return cacheValue, nil
		}
	}

	value, err := getValue(key)
	if err != nil {
		slog.Warn("Used fallback", "ticker", key, "cacheValue", cacheValue)
		return cacheValue, nil
	}

	dbSave(key, Cache[V]{
		Value:   value,
		Created: time.Now(),
	})

	return value, nil
}
