# async

A ES7-style async/await implementation lib in Golang.

## API

- `async.AsyncTask`

Async work unit created by `async.Async`. 

```go
type AsyncTask[V any] struct {
	Value   V  // execution result
	Err     any // execution possible error
	State   AsyncState // execution state
}
```

`async.AsyncTask.State` has 3 values: `async.PENDING`(`async.AsyncTask` is doing), `async.FULFILLED`(`async.AsyncTask` is done and result returns), `async.REJECTED`(`async.AsyncTask` is done but error happens). 

- `async.Async`

Convert a synchronous `func` to be asynchronous.

- `async.Await`

Wait `async.AsyncTask` to be done.

- `async.All`

Asynchronously return values(when all `func`s are finished) or error(once one of `func`s throws error).

- `async.Race`

Asynchronously return value or error of the fisrt finished `func`.

