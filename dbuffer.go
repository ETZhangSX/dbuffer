package dbuffer

import (
	"sync"
	"sync/atomic"
	"time"
)

// DBuffer is the interface for double-buffer that support data hot update.
type DBuffer[T any] interface {
	// Data returns the stored data.
	Data() T
	// DataWithDone returns the stored data and a DoneFunc while increasing
	// the ref counter by one. The ref counter decrease when DoneFunc is called.
	DataWithDone() (data T, done DoneFunc)
	// Load is used to load new data manually. It will be blocked if the writing
	// buffer ref counter is larger than 0.
	Load()
}

var _ DBuffer[any] = (*dBuffer[any])(nil)

// Alloc is the func type for allocate the object.
type Alloc[T any] func() T

// New creates a DBuffer instance.
func New[T any](loader Loader[T], alloc Alloc[T], opts ...Option) DBuffer[T] {
	buf := &dBuffer[T]{
		loader: loader,
		refs:   make([]sync.WaitGroup, 2),
		opts:   &Options{},
	}
	buf.data = append(buf.data, alloc(), alloc())
	buf.opts.init()
	for _, o := range opts {
		o(buf.opts)
	}
	buf.watch()
	return buf
}

// Implementation of DBuffer.
type dBuffer[T any] struct {
	loader Loader[T]
	opts   *Options

	data []T
	refs []sync.WaitGroup
	idx  int32
}

func (d *dBuffer[T]) watch() {
	d.Load()
	if d.opts.interval < 0 {
		return
	}
	go func() {
		for {
			time.Sleep(d.opts.interval)
			d.Load()
		}
	}()
}

func (d *dBuffer[T]) Load() {
	i := 1 - atomic.LoadInt32(&d.idx)
	// if there are still some refs, wait for all done.
	d.refs[i].Wait()
	if ok, _ := d.loader.Load(&d.data[i]); ok {
		atomic.StoreInt32(&d.idx, i)
	}
}

func (d *dBuffer[T]) Data() T {
	i := atomic.LoadInt32(&d.idx)
	return d.data[i]
}

// DoneFunc decrements the ref counter of data by one.
type DoneFunc func()

func (d *dBuffer[T]) DataWithDone() (data T, done DoneFunc) {
	i := atomic.LoadInt32(&d.idx)
	r := &d.refs[i]
	r.Add(1)
	return d.data[i], func() {
		r.Done()
	}
}
