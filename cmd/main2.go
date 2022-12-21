package main

import (
	"fmt"
	"time"
)

func main2() {

	d := New(2)

	d.Debounce(SayHello, 1)
	d.Debounce(SayHello, 2)

	time.Sleep(2 * time.Second)

	d.Debounce(SayHello, 3)
	d.Debounce(SayHello, 4)

}

func SayHello(i int) {
	fmt.Printf("Hello %d\n", i)
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

func (d *Debounce) Debounce(f func(int), i int) {

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
