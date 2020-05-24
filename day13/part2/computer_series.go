package main

import (
	"sync"

	"github.com/pkg/errors"
)

type seriesComputer struct {
	wg         sync.WaitGroup
	comps      []*computer
	inputChans []chan int64
	outputChan chan int64
}

func newSeriesComputer(inputText string, labels ...string) (*seriesComputer, error) {
	template, err := newChannelComputer(inputText, nil, nil)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert input to int code memory")
	}

	series := &seriesComputer{}

	connectChan := make(chan int64, 1)
	for _, label := range labels {
		series.inputChans = append(series.inputChans, connectChan)
		next := make(chan int64, 1)

		comp := template.clone()
		comp.inChan = connectChan
		comp.outChan = next
		comp.label = label
		comp.disableLog = true
		series.comps = append(series.comps, comp)

		connectChan = next
	}
	series.outputChan = connectChan
	return series, nil
}

func newFeedbackComputer(inputText string, labels ...string) (*seriesComputer, error) {
	series, err := newSeriesComputer(inputText, labels...)
	if err != nil {
		return nil, errors.Wrap(err, "unable to convert input to int code memory")
	}

	// connect output back to first input
	series.comps[0].inChan = series.outputChan
	series.inputChans[0] = series.outputChan
	return series, nil
}

// runAsync starts computers running which will block until inputs are met
func (s *seriesComputer) runAsync() {
	for _, comp := range s.comps {
		cToRun := comp
		s.wg.Add(1)
		go func() {
			defer s.wg.Done()
			err := cToRun.run()
			if err != nil {
				cToRun.logf("error running computer: %s\n", err.Error())
			}
		}()
	}
}

func (s *seriesComputer) waitForCompletion() {
	s.wg.Wait()
}

func (s *seriesComputer) loadPhases(phases []int) error {
	if len(phases) != len(s.comps) {
		return errors.Errorf("incorrect number of phases provided (got %d, want %d)", len(phases), len(s.comps))
	}

	for i, phase := range phases {
		s.inputChans[i] <- int64(phase)
	}
	return nil
}

func (s *seriesComputer) input(arg int64) {
	s.inputChans[0] <- arg
}

func (s *seriesComputer) output() int64 {
	return <-s.outputChan
}
