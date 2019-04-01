package bufferexecutor

import (
	"log"
	"time"
)

type bufferExecutor struct {
	bufSize  int
	interval int
	work     func([]interface{})

	cached []interface{}
	dataCh chan interface{}
	ticker *time.Ticker
	signal chan struct{}
	stop   chan struct{}
}

type BufferExecutor interface {
	Put(data interface{})
	Signal()
	Close() error
	Size() int
}

const MaxSignal = 100

func New(bufSize int, interval int, work func([]interface{})) BufferExecutor {
	if work == nil {
		panic("work must not be nil")
	}

	if bufSize < 1 {
		panic("bufSize must greater than zero")
	}

	if interval < 1 {
		panic("interval must greater than zero")
	}

	if interval < 100 {
		interval = 100
	}

	ps := bufferExecutor{
		bufSize:  bufSize,
		interval: interval,
		work:     recoverF(work),
	}

	ps.cached = make([]interface{}, 0, bufSize)
	ps.dataCh = make(chan interface{}, bufSize)
	ps.ticker = time.NewTicker(time.Duration(int64(ps.interval) * int64(time.Millisecond)))
	ps.signal = make(chan struct{}, MaxSignal)
	ps.stop = make(chan struct{})

	go func() {
		for {
			select {
			case item := <-ps.dataCh:
				ps.cached = append(ps.cached, item)
				if len(ps.cached) == cap(ps.cached) {
					ps.flush(ps.cached)
				}

				break
			case <-ps.ticker.C:
				ps.flush(ps.cached)
				break
			case <-ps.signal:
				ps.flush(ps.cached)
				break
			case <-ps.stop:
				ps.clear()

				for item := range ps.dataCh {
					if len(ps.cached) == cap(ps.cached) {
						ps.flush(ps.cached)
					}
					ps.cached = append(ps.cached, item)

				}

				ps.flush(ps.cached)

				return

			}
		}
	}()

	return &ps
}

func (ps *bufferExecutor) flush(data []interface{}) {
	ps.work(data)
	ps.cached = ps.cached[:0]
}

func recoverF(f func([]interface{})) func([]interface{}) {
	return func(data []interface{}) {
		defer func() {
			if r := recover(); r != nil {
				log.Println(r)
			}
		}()
		f(data)
	}
}

func (ps *bufferExecutor) Signal() {
	ps.signal <- struct{}{}
}

func (ps *bufferExecutor) Put(data interface{}) {
	ps.dataCh <- data
}

func (ps *bufferExecutor) Close() error {
	ps.stop <- struct{}{}

	return nil
}

func (ps *bufferExecutor) clear() {
	close(ps.dataCh)
	ps.ticker.Stop()
	close(ps.signal)
	close(ps.stop)
}

func (ps *bufferExecutor) Size() int {
	return len(ps.cached)
}
