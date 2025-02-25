package repository

type MemoryRepository interface {
	Set(key string, value int64) error
	Get(key string) (int64, bool)
	Exists(key string) bool
}
