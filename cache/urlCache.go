package cache

import (
	"errors"
	"time"

	"github.com/EfoJensen/go-rentrospect/upload"
	"github.com/dgraph-io/ristretto"
)

type UrlCache struct {
	ttl     time.Duration
	storage *upload.Storage
	cache   *ristretto.Cache
}

func NewUrlCache(storage *upload.Storage) (*UrlCache, error) {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,
		MaxCost:     1 << 30,
		BufferItems: 64,
	})

	if err != nil {
		return nil, err
	}

	return &UrlCache{
		cache:   cache,
		storage: storage,
		ttl:     23 * time.Hour,
	}, err
}

func (c *UrlCache) GetURL(objectName string) (string, error) {
	if url, found := c.cache.Get(objectName); found {
		return url.(string), nil
	}

	url, err := c.storage.GetTempUrl(objectName, 24)

	if err != nil {
		return "", err
	}

	if success := c.cache.SetWithTTL(objectName, url, 1, c.ttl); !success {
		return "", errors.New("high write contention; internal buffer full")
	}

	return url, nil
}

func (c *UrlCache) GetUrlBatch(objectNames []string) (map[string]string, error) {
	var missing []string
	results := make(map[string]string)

	for _, name := range objectNames {
		if url, found := c.cache.Get(name); found {
			results[name] = url.(string)
		} else {
			missing = append(missing, name)
		}
	}

	for _, name := range missing {
		url, err := c.storage.GetTempUrl(name, 24)

		if err != nil {
			return nil, err
		}
		c.cache.SetWithTTL(name, url, 1, 23)

		results[name] = url
	}

	return results, nil
}