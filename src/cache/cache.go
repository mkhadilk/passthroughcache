package cache

type handler[T any] interface {
	Get(id string) (T, error)
	Refresh(id string) error
}
