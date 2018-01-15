package tcp

import (
	"testing"

	"github.com/reactivex/rxgo/handlers"
	"github.com/reactivex/rxgo/observable"
	"github.com/reactivex/rxgo/observer"
	"net/http"
	"runtime"
	"sync"
)

func init() {
	go http.ListenAndServe("127.0.0.1:3000", nil)
}
func BenchmarkChan(b *testing.B) {
	runtime.GOMAXPROCS(3)
	num := 0
	wg := sync.WaitGroup{}
	for i := 0; i < b.N; i++ {
		wg.Add(1)
		go func() {
			_, err := http.DefaultClient.Get("http://127.0.0.1:3000")
			if err != nil {
				b.Error(err)
			}
			num++
			wg.Done()

		}()
	}
	wg.Wait()
	b.Log("Done:", num)
}

func BenchmarkRX(b *testing.B) {
	runtime.GOMAXPROCS(3)
	num := 0
	wg := sync.WaitGroup{}
	onNext := handlers.NextFunc(func(item interface{}) {
		_, err := http.DefaultClient.Get("http://127.0.0.1:3000")
		if err != nil {
			b.Error(err)
		}
	})

	onDone := handlers.DoneFunc(func() {
		num++
	})

	watcher := observer.New(onNext, onDone)
	b.ResetTimer()
	// Create an `Observable` from a single item and subscribe to the observer.
	for i := 0; i < b.N; i++ {
		go func() {
			wg.Add(1)
			sub := observable.Just(1).Subscribe(watcher)
			<-sub
			num++
			wg.Done()

		}()
	}
	wg.Wait()
	b.Log("Done:", num)
}
