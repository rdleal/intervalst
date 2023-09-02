# Interval Search Tree

[![Go Reference](https://pkg.go.dev/badge/github.com/rdleal/intervalst/interval.svg)](https://pkg.go.dev/github.com/rdleal/intervalst/interval)
[![Go Report Card](https://goreportcard.com/badge/github.com/rdleal/intervalst)](https://goreportcard.com/report/github.com/rdleal/intervalst)
[![codecov](https://codecov.io/gh/rdleal/intervalst/branch/main/graph/badge.svg?token=BMJKC8DT9U)](https://codecov.io/gh/rdleal/intervalst)

Package [interval](./interval) provides a generic interval tree implementation.

An interval tree is a data structure useful for storing values associated with intervals,
and efficiently search those values based on intervals that overlap with any given interval.
This generic implementation uses a self-balancing binary search tree algorithm, so searching
for any intersection has a worst-case time-complexity guarantee of <= 2 log N, where N is the number of items in the tree.

For more on interval trees, [see](https://en.wikipedia.org/wiki/Interval_tree)

## Usage

Importing the package:
```go
import "github.com/rdleal/intervalst/interval"
```
Creating a tree with `time.Time` as interval key type and `string` as value type:
```go
cmpFn := func(t1, t2 time.Time) int {
        switch{
        case t1.After(t2): return 1
        case t1.Before(t2): return -1
        default: return 0
        }
}
st := interval.NewSearchTree[string](cmpFn)
```

Upserting a value:
```go
start := time.Now()
end := start.Add(2*time.Hour)
err := st.Insert(start, end, "event 1")
if err != nil {
        // error handling...
}
```
Searching for any intersection:
```go
start := time.Now()
end := start.Add(2*time.Hour)
val, ok := st.AnyIntersection(start, end)
if !ok {
        // no intersection found for start and end...
}
```

Deleting an interval from the tree:
```go
start := time.Now()
end := start.Add(2*time.Hour)
err := st.Delete(start, end)
if err != nil {
        // error handling...
}
```

For more operations, check out the [GoDoc page](https://pkg.go.dev/github.com/rdleal/intervalst/interval).

## Testing

Running unit tests:
```sh
$ go test -v -cover -race ./...
```

## Benchmarks

Running benchmarks:
```sh
$ cd interval && go test -run='^$' -bench=. -benchtime=20s -benchmem
```

## License

[MIT](./LICENSE)
