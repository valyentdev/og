package storage

type Driver interface {
	Upload(key string, data []byte) error
	Get(key string) ([]byte, error)
}
