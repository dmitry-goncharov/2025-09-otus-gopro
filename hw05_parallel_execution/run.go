package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	stop := make(chan struct{}, n+1)

	jobs := producer(stop, tasks)

	errs := make(chan error)
	defer close(errs)

	wg := &sync.WaitGroup{}

	for range n {
		wg.Add(1)
		go worker(wg, jobs, errs, stop)
	}

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

func producer(stop <-chan struct{}, tasks []Task) <-chan Task {
	jobs := make(chan Task)

	go func() {
		defer close(jobs)

		for _, task := range tasks {
			select {
			case <-stop:
				return
			case jobs <- task:
			}
		}
	}()

	return jobs
}

func worker(wg *sync.WaitGroup, tasks <-chan Task, errs chan<- error, stop <-chan struct{}) {
	defer wg.Done()

	for {
		select {
		case <-stop:
			return
		case task, ok := <-tasks:
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
