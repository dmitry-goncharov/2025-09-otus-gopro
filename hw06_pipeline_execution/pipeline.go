package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, stage := range stages {
		ch = stage(take(done, ch))
	}
	return ch
}

func take(done In, in In) Out {
	out := make(Bi)
	go func() {
		defer func() {
			close(out)
			for v := range in {
				_ = v
			}
		}()
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
				case out <- v:
				}
			}
		}
	}()
	return out
}
