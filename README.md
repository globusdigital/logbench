<a href="https://goreportcard.com/report/github.com/globusdigital/logbench">
  <img src="https://goreportcard.com/badge/github.com/globusdigital/logbench" alt="GoReportCard">
</a>
<a href="https://travis-ci.org/globusdigital/logbench">
    <img src="https://travis-ci.org/globusdigital/logbench.svg?branch=master" alt="Travis CI: build status">
</a>
<a href='https://coveralls.io/github/globusdigital/logbench'>
    <img src='https://coveralls.io/repos/github/globusdigital/logbench/badge.svg' alt='Coverage Status' />
</a>

<h2>
  <span>Go - Structured Logging Benchmark</span>
  <br>
  <sub>by <a href="globus.ch">Magazine zum Globus</a></sub>
</h2>

A structured JSON logging performance benchmark providing realistic performance metrics for the latest versions of:
- [sirupsen/logrus](https://github.com/sirupsen/logrus)
- [uber/zap](https://github.com/uber-go/zap)
- [rs/zerolog](https://github.com/rs/zerolog)

Performance is measured by the following main criteria:
- `total alloc` - Total size of allocated memory.
- `num-gc` - Total number of GC cycles.
- `mallocs` - Total number of allocated heap objects.
- `total pause` - Total duration of GC pauses.
- Average and total time of execution (by operation).

<br>

![Benchmark Results](https://github.com/globusdigital/logbench/blob/master/results_graphs.png?raw=true)
_i7-8569U @ 2.80GHz_

## Getting started

1. Install the benchmark:
```
go get github.com/globusdigital/logbench
```

2. Run it:
```
logbench -w 8 -t 1_000_000 -o_all -l zerolog
```

- `-w <num>`: defines the number of concurrently writing goroutines.
- `-l <logger>`: enables a logger.
You can enable multiple loggers by specifying multiple flags: `-l zerolog -l zap -l logrus`.
- `-o <operation>`: enables an operation.
You can enable multiple operations by specifying multiple flags: `-o info -o error -o info_with_3`.
- `-t <num>`: defines the number of logs to be written for each operation.
- `-o_all`: enables all operations ignoring all specified `-o` flags
- `-memprof <path>`: specifies the output file path for the memory profile (disabled when not set)
- `-mi <duration>`: memory inspection interval

## How-to
### Adding a new logger to the benchmark
- 1. Define the logger in a sub-package.
- 2. Provide a `Setup() benchmark.Setup` function in your logger's sub-package.
- 3. Implement all benchmark operations:
  - `FnInfo func(msg string)`
  - `FnInfoFmt func(msg string, data int)`
  - `FnError func(msg string)`
  - `FnInfoWithErrorStack func(msg string, err error)`
  - `FnInfoWith3 func(msg string, fields *benchmark.Fields3)`
  - `FnInfoWith10 func(msg string, fields *benchmark.Fields10)`
  - `FnInfoWith10Exist func(msg string)`
- 4. Add your setup to [`setups`](https://github.com/globusdigital/logbench/blob/eff659cfb1eb06b1d139db6735b2b2ce6944632c/main.go#L21).
- 5. Run the tests with `go test -v -race ./...` and make sure everything's working.
