package onepiece

import "sync"

type TypeFactory func() interface{}

type TypeProvider struct {
	Types map[string]TypeFactory
	mu    sync.RWMutex
}

func (t *TypeProvider) Register(name string, factory TypeFactory) *TypeProvider {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.Types[name] = factory
	return t
}

func (t *TypeProvider) Get(name string) (TypeFactory, bool) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	factory, ok := t.Types[name]
	return factory, ok
}

func NewTypeProvider() *TypeProvider {
	return &TypeProvider{
		Types: make(map[string]TypeFactory),
	}
}

func GenericFactory[Message any]() interface{} {
	var m Message
	return &m
}

func FetchGeneric[Message any](tp *TypeProvider, name string) (*Message, error) {
	if factory, ok := tp.Get(name); ok {
		payload := factory()

		switch payload.(type) {
		case Message:
			return payload.(*Message), nil
		default:
			return nil, ErrUnknownMessage
		}
	}

	return nil, ErrUnknownMessage
}
