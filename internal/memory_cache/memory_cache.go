package memorycache

import (
	"sync"
	"time"
)

// cacheItem теперь хранит список значений и время истечения
type cacheItem struct {
	values    []string
	expiresAt time.Time
}

type MemoryCache struct {
	cache map[string]cacheItem
	mu    sync.Mutex
}

func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		cache: make(map[string]cacheItem),
	}
}

// получение записи из кэша
func (m *MemoryCache) Get(key string) ([]string, bool) {
	m.mu.Lock()
	defer m.mu.Unlock()

	item, found := m.cache[key]
	if !found {
		return nil, false
	}

	// проверяем, истек ли срок действия записи
	if time.Now().After(item.expiresAt) {
		delete(m.cache, key) // удаляем устаревшую запись
		return nil, false
	}

	return item.values, true
}

// внесение записи в кэш с TTL
func (m *MemoryCache) Set(key string, values []string, ttl time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expiresAt := time.Now().Add(ttl)
	m.cache[key] = cacheItem{
		values:    values,
		expiresAt: expiresAt,
	}
}

// проверка наличия данных
func (m *MemoryCache) IsExists(key string) bool {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, found := m.cache[key]
	return found
}

// удаление записи из кэша
func (m *MemoryCache) Delete(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.cache, key)
}

// очистка устаревших записей
func (m *MemoryCache) Cleanup() {
	m.mu.Lock()
	defer m.mu.Unlock()

	for key, item := range m.cache {
		if time.Now().After(item.expiresAt) {
			delete(m.cache, key)
		}
	}
}
