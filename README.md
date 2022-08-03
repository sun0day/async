# async

A ES7-style async/await implementation lib in Golang.

## Install

```shell
$ go get github.com/sun0day/async
```

## API

### `async.AsyncTask`

Async work unit created by `async.Async`. 

```go
type AsyncTask[V any] struct {
  Value   V  // execution result
  Err     any // execution possible error
  State   AsyncState // execution state
}
```

`async.AsyncTask.State` has 3 values: `async.PENDING`(`async.AsyncTask` is doing), `async.FULFILLED`(`async.AsyncTask` is done and result returns), `async.REJECTED`(`async.AsyncTask` is done but error happens). 

### `async.Async`

Convert a synchronous `func` to be asynchronous.

**Type definition**

```go
func Async[V any](f func() V) func() *AsyncTask[V] {...}
```

**Usage**

```go
package main

import (
  "fmt"
  "time"
  "github.com/sun0day/async"
)

func main() {
  a := 1
  b := 2
  f := func() int {
    c := a + b
    fmt.Printf("async result=%d\n", c)
    return c
  }

  af := async.Async[int](f)

  fmt.Printf("sync start, goroutine=%d\n", runtime.NumGoroutine())
  af()
  fmt.Printf("sync end, goroutine=%d\n", runtime.NumGoroutine())
  time.Sleep(1 * time.Second)
  fmt.Printf("async end, goroutine=%d\n", runtime.NumGoroutine())
}

/* stdout
sync start, goroutine=1
sync end, goroutine=2
async result=3
async end, goroutine=1
*/
```

### `async.Await`

Wait `async.AsyncTask` to be done.

**Type definition**

```go
func Await[V any](t *AsyncTask[V]) (V, any) {...}
```

**Usage**

```go
package main

import (
  "fmt"
  "runtime"
  "github.com/sun0day/async"
)

func main() {
  a := 1
  b := 2
  f1 := func() int {
    c := a + b
    return c
  }

  f2 := func() int {
    panic("f() error")
  }

  af1 := async.Async[int](f1)
  af2 := async.Async[int](f2)

  fmt.Printf("sync start, goroutine=%d\n", runtime.NumGoroutine())
  value, _ := async.Await[int](af1())
  _, err := async.Await[int](af2())
  fmt.Printf("af1 result=%d\n", value)
  fmt.Printf("af2 error=%s\n", err)
  fmt.Printf("sync end, goroutine=%d\n", runtime.NumGoroutine())
}
/* stdout
sync start, goroutine=1
af1 result=3
af2 error=f() error
sync end, goroutine=1
*/
```

### `async.All`

Asynchronously return values(when all `func`s are finished) or error(once one of `func`s throws error).

**Type definition**

```go
func All[V any](fs []func() V) *AsyncTask[[]V] {...}
```

**Usage**

```go
package main

import (
  "fmt"
  "runtime"
  "time"
  "github.com/sun0day/async"
)

func main() {
  a := 1
  b := 2
  f1 := func() int {
    time.Sleep(1 * time.Second)
    c := a + b
    return c
  }

  f2 := func() int {
    time.Sleep(2 * time.Second)
    c := a * b
    return c
  }

  f3 := func() int {
    panic("f3 error")
  }

  fmt.Printf("sync start, goroutine=%d\n", runtime.NumGoroutine())
  values, _ := async.Await[[]int](async.All[int]([]func() int{f1, f2}))
  _, err := async.Await[[]int](async.All[int]([]func() int{f1, f2, f3}))
  fmt.Printf("all result=%v\n", values)
  fmt.Printf("all error=%s\n", err)
  fmt.Printf("sync end, goroutine=%d\n", runtime.NumGoroutine())
}
/* stdout
sync start, goroutine=1
all result=[3 2]
all error=f3 error
sync end, goroutine=1
*/
```

### `async.Race`

Asynchronously return value or error of the fisrt finished `func`.

**Type definition**

```go
func Race[V any](fs []func() V) *AsyncTask[V] {...}
```

**Usage**

```go
package main

import (
  "fmt"
  "runtime"
  "time"
  "github.com/sun0day/async"
)

func main() {
  a := 1
  b := 2
  f1 := func() int {
    time.Sleep(1 * time.Second)
    c := a + b
    return c
  }

  f2 := func() int {
    time.Sleep(2 * time.Second)
    c := a * b
    return c
  }

  fmt.Printf("sync start, goroutine=%d\n", runtime.NumGoroutine())
  value, _ := async.Await[int](async.Race[int]([]func() int{f1, f2}))
  fmt.Printf("race result=%d\n", value)
  fmt.Printf("sync end, goroutine=%d\n", runtime.NumGoroutine())
}
/* stdout
sync start, goroutine=1
race result=3
sync end, goroutine=2
*/
```

