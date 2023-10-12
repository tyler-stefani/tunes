package cache

type Cache interface {
	Get(key string) string
	Put(key string, value string)
}
