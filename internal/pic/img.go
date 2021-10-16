package pic

import (
	"fmt"
	"io"
	"sync"
)

type img struct {
	pool  *sync.Pool
	alive bool
	mutex sync.Mutex
	buf   []byte
}

func (i *img) Recycle() {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if !i.alive {
		return
	}

	i.alive = false
	i.pool.Put(i)
}

func (i *img) WriteTo(w io.Writer) (n int64, err error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()

	if !i.alive {
		return 0, fmt.Errorf("image is already recycled")
	}

	nb, err := w.Write(i.buf)

	return int64(nb), err
}
