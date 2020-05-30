package main

import "sync"

type inOutComputer struct {
	wg   sync.WaitGroup
	comp *computer
	in   chan<- int64
	out  <-chan int64
}

func newInOutComputer(inputText string) (*inOutComputer, error) {
	in := make(chan int64)
	out := make(chan int64)
	c, err := newChannelComputer(inputText, in, out)
	if err != nil {
		return nil, err
	}
	c.disableLog = true
	c.disableOutLog = true

	result := &inOutComputer{
		comp: c,
		in:   in,
		out:  out,
	}

	result.wg.Add(1)
	go func() {
		c.run()
		result.wg.Done()
	}()
	return result, nil
}

func (i *inOutComputer) Input(val int) int {
	i.in <- int64(val)
	return int(<-i.out)
}
