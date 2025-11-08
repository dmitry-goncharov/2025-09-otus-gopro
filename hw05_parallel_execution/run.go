package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	jobs := make(chan Task, len(tasks))

	var errsLen int
	if m > 0 {
		errsLen = 0
	} else {
		errsLen = len(tasks)
	}
	errs := make(chan error, errsLen)
	defer close(errs)

	stop := make(chan struct{}, n)

	wg := &sync.WaitGroup{}

	for range n {
		wg.Add(1)
		go worker(wg, jobs, errs, stop)
	}

	go func() {
		defer close(jobs)

		for _, task := range tasks {
			jobs <- task
		}
	}()

	var res error
	go func() {
		defer close(stop)

		errCounter := 0
		for range errs {
			if m > 0 {
				errCounter++
				if errCounter == m {
					res = ErrErrorsLimitExceeded
					for range n {
						stop <- struct{}{}
					}
				}
			}
		}
	}()

	wg.Wait()

	return res
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errs chan<- error, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		default:
			task, ok := <-tasks
			if ok {
				err := task()
				if err != nil {
					errs <- err
				}
			} else {
				return
			}
		}
	}
}
