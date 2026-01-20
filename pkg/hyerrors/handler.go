package hyerrors

import "sync"

type Handler interface {
	Handle(err *Error)
}

type HandlerFunc func(err *Error)

func (f HandlerFunc) Handle(err *Error) {
	f(err)
}

type Registry struct {
	mu       sync.RWMutex
	handlers []Handler
}

var global = &Registry{}

func RegisterHandler(h Handler) {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.handlers = append(global.handlers, h)
}

func RegisterHandlerFunc(fn func(err *Error)) {
	RegisterHandler(HandlerFunc(fn))
}

func ClearHandlers() {
	global.mu.Lock()
	defer global.mu.Unlock()
	global.handlers = nil
}

func Handle(err *Error) {
	if err == nil {
		return
	}

	global.mu.RLock()
	handlers := make([]Handler, len(global.handlers))
	copy(handlers, global.handlers)
	global.mu.RUnlock()

	for _, h := range handlers {
		h.Handle(err)
	}
}

func Report(err *Error) {
	Handle(err)
}
