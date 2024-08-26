package bufferpool

import "sync"

var (
	_pool = sync.Pool{
		New: func() any {
			return &Buffer{}
		},
	}
)

func Get() *Buffer {
	buf := _pool.Get().(*Buffer)
	buf.Reset()

	return buf
}
