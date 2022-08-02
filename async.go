package async

type AsyncState int

const (
	PENDING AsyncState = iota
	FULFILLED
	REJECTED
)

type AsyncTask[V any] struct {
	valueCh chan V
	errCh   chan any
	value   V
	err     any
	state   AsyncState
}

func (t *AsyncTask[V]) Value() V {
	return t.value
}

func Async[V any](f func() V) func() *AsyncTask[V] {
	exec := func() *AsyncTask[V] {
		task := &AsyncTask[V]{
			valueCh: make(chan V, 1),
			errCh:   make(chan any, 1),
			state:   PENDING,
		}

		go func() {
			var result V
			defer func() {
				if e := recover(); e == nil {
					task.state = FULFILLED
					task.value = result
					task.valueCh <- result
				} else {
					task.state = REJECTED
					task.err = e
					task.errCh <- e
				}
				close(task.valueCh)
				close(task.errCh)
			}()

			result = f()
		}()

		return task
	}

	return exec
}

func Await[V any](t *AsyncTask[V]) (V, any) {
	var value V
	var err any

	switch t.state {
	default:
		select {
		case err = <-t.errCh:
			return value, err
		case value = <-t.valueCh:
			return value, err
		}
	case FULFILLED:
		return t.value, err
	case REJECTED:
		return value, t.err
	}

}

func All[V any](fs []func() V) *AsyncTask[[]V] {
	return Async[[]V](func() []V {
		count := len(fs)
		tasks := make([]*AsyncTask[V], count)
		values := make([]V, count)
		for i, f := range fs {
			tasks[i] = Async[V](f)()
		}

		for i, t := range tasks {
			value, err := Await[V](t)
			values[i] = value
			if err != nil {
				panic(err)
			}
		}

		return values
	})()
}

func Race[V any](fs []func() V) *AsyncTask[V] {
	return Async[V](func() V {
		count := len(fs)
		valueCh := make(chan V, count)
		errCh := make(chan any, count)
		for i := range fs {
			go func(i int) {
				var value V
				f := fs[i]

				defer func() {
					if e := recover(); e == nil {
						valueCh <- value
					} else {
						errCh <- e
					}
				}()

				value = f()
			}(i)
		}

		select {
		case value := <-valueCh:
			return value
		case err := <-errCh:
			panic(err)
		}
	})()
}
