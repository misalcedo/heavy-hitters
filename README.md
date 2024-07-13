## Heavy Hitters
Implementation of the [SpaceSaving](https://www.cs.ucsb.edu/sites/default/files/documents/2005-23.pdf) algorithm.

### Simulation
```console
go run examples/simulation.go
```

### Tests
```console
go test ./...
```

### Benchmark
```console
go test -bench=. -run=^$ ./...
```

### Example
The following command will:
```console
go run ./... $FILE
```

1. Read the contents of the file at path `$FILE` into memory.
2. Split the file by whitespace.
3. Approximate the frequent and top-6 elements in the file using the SpaceSaving algorithm.

Each whitespace-separated entry in the file is one element.
For example, the following while would result in a stream of `[]string{"1", "2", "3", "4", "5"}`: 

```text
1
2
3
4
5
```
