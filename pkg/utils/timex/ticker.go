package timex

import "time"

func FnTicker(d time.Duration, quit <-chan struct{}, fn func()) {
	tick := time.NewTicker(d)

	go func() {
		for {
			select {
			case <-tick.C:
				fn()
			case <-quit:
				return
			}
		}
	}()
}
