package main

import (
	"fmt"
	"time"
)

func main() {

	{
		//d := New(2)
		//d.Debounce(SayHello, 1)
	}

	dur := time.Duration(5) * time.Second

	in := make(chan int) // there's input data stream
	ex := make(chan bool)
	out := debounce3(in, dur, ex) // there's out data set

	// Per second tact's imitation
	go func(in chan int) {
		{
			//for _, e := range [1]int{1} { // out: 1
			//	in <- e
			//	time.Sleep(4 * time.Second)
			//}
			//for _, e := range [5]int{1, 2, 3, 4, 5} { // out: 1
			//	in <- e
			//	time.Sleep(1 * time.Second)
			//}
			//for _, e := range [6]int{1, 2, 3, 4, 5, 6} { // out: 1, 6
			//	in <- e
			//	time.Sleep(1 * time.Second)
			//}
		}
		for _, e := range [8]int{1, 2, 3, 4, 5, 6, 7, 8} { // out: 1, 6
			in <- e
			time.Sleep(1 * time.Second)
		}
		close(in)
		ex <- true
		close(ex)
	}(in)

	for i := range out { // data out
		fmt.Printf("out: %d\n", i)
	}
}

type Debounce struct {
	t        *time.Ticker
	debounce time.Duration
	past     int
	first    bool
}

func New(debounce int) *Debounce {
	return &Debounce{
		t:        time.NewTicker(1 * time.Second),
		debounce: time.Duration(debounce),
		past:     0,
		first:    true,
	}
}

func (d *Debounce) debounce0(f func(int), i int) {

	go func() {
		for {
			select {
			case <-d.t.C:
				d.past++
			}
		}
	}()

	now := time.Now()
	p := now.Add(-time.Duration(d.past))
	if now.Sub(p) >= d.debounce || d.first {
		d.past = 0
		d.first = false
		f(i)
	}

}

// debounce
func debounce1(in chan int, interval time.Duration) chan int {
	out := make(chan int)
	now := time.Now()
	count := 0
	go func(out chan int, now time.Time, interval time.Duration) {
		for i := range in {
			passed := time.Now().Sub(now)
			// fmt.Println(passed)
			if passed >= interval || count == 0 {
				out <- i
				now = time.Now()
				count = 0
			}
			count++
		}
		close(out)
	}(out, now, interval)
	return out
}

func debounce2(in chan int, dur time.Duration, ex chan bool) chan int {

	out := make(chan int)
	go func(in chan int, out chan int, dur time.Duration) {

		timed := false
		loop := true
		exited := false
		for loop {
			select {
			case i := <-in:
				if !timed {
					timed = true
					go func(out chan int, i int, dur time.Duration) {
						time.Sleep(dur)
						timed = false
						out <- i
						if exited {
							loop = false
						}
					}(out, i, dur)
				}
			case <-ex:
				exited = true
			}
		}

		close(out)

	}(in, out, dur)

	return out
}

func debounce3(in chan int, dur time.Duration, ex chan bool) chan int {

	out := make(chan int)
	go func(in chan int, out chan int, dur time.Duration) {

		loop := true
		exited := false
		timed := false

		var swap int
		t := time.NewTimer(dur)
		for loop {
			select {
			case i := <-in:
				if !timed {
					timed = true
					t.Reset(dur)
					swap = i
				}
			case <-ex:
				exited = true
			case <-t.C:
				out <- swap
				swap = 0
				timed = false
				if exited {
					loop = false
				}
			}
		}

		close(out)

	}(in, out, dur)

	return out
}
