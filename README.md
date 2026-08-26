# Go Tutorial

This repository is a beginner-friendly Go language learning project. It contains small, focused Go programs that introduce core concepts step by step — from language basics all the way through concurrency, packages, web development, databases, and testing/production practices.

## Purpose

The goal of this repo is to teach the fundamentals of Go programming through simple examples and practical exercises. Each file builds on the previous one and covers essential topics such as:

- Variables, constants, and operators
- Conditions and loops
- Arrays, slices, and maps
- Functions, closures, and recursion
- Pointers, structs, methods, and interfaces
- Enums and generics
- Error handling, panic, and recover
- File handling, JSON, and CSV
- Goroutines, channels, and concurrency patterns (mutexes, select, worker pools, fan-in/fan-out, context)
- Packages and modules
- Web development (HTTP servers/clients, REST APIs, middleware, auth)
- Databases (SQL basics, CRUD, transactions, connection pooling, migrations)
- Testing and production readiness (unit/table-driven tests, benchmarks, profiling, logging, deployment)

## Files

### Language basics

- [1_hello.go](1_hello.go) — Hello World example
- [2_value_variable.go](2_value_variable.go) — Variables and values
- [3_constants.go](3_constants.go) — Constants
- [4_operators.go](4_operators.go) — Arithmetic and comparison operators
- [5_conditions.go](5_conditions.go) — If/else conditions
- [6_loops.go](6_loops.go) — Loop structures
- [7_arrays.go](7_arrays.go) — Arrays
- [8_slices.go](8_slices.go) — Slices
- [9_maps.go](9_maps.go) — Maps
- [10_range.go](10_range.go) — Range keyword
- [11_functions.go](11_functions.go) — Function definitions and usage
- [12_closures.go](12_closures.go) — Closures
- [13_recursion.go](13_recursion.go) — Recursion

### Types and structures

- [14_pointers.go](14_pointers.go) — Pointers
- [15_structs.go](15_structs.go) — Structs
- [16_methods_receiver.go](16_methods_receiver.go) — Methods and receivers
- [17_interfaces.go](17_interfaces.go) — Interfaces, type assertion, and type switch
- [18_enums.go](18_enums.go) — Enums with `iota`
- [19_generics.go](19_generics.go) — Generics

### Error handling and I/O

- [20_error_handling.go](20_error_handling.go) — Error handling
- [21_defer_panic_recover.go](21_defer_panic_recover.go) — Defer, panic, and recover
- [22_file_handling.go](22_file_handling.go) — File handling
- [23_json_csv.go](23_json_csv.go) — JSON and CSV encoding/decoding

### Concurrency

- [24_goroutines.go](24_goroutines.go) — Goroutines
- [25_channels.go](25_channels.go) — Channels
- [26_buffered_channels.go](26_buffered_channels.go) — Buffered channels
- [27_directional_channels.go](27_directional_channels.go) — Directional channels
- [28_select.go](28_select.go) — Select statement
- [29_mutex.go](29_mutex.go) — Mutex
- [30_rwmutex.go](30_rwmutex.go) — RWMutex
- [31_race_conditions.go](31_race_conditions.go) — Race conditions
- [32_deadlocks.go](32_deadlocks.go) — Deadlocks
- [33_worker_pools.go](33_worker_pools.go) — Worker pools
- [34_fan_in_fan_out.go](34_fan_in_fan_out.go) — Fan-in/fan-out
- [35_context_cancellation.go](35_context_cancellation.go) — Context cancellation

### Packages, modules, and projects

- [36_packages_modules_project/](36_packages_modules_project/) — A small multi-package project (`cmd/app`, `internal/mathutil`, `stringutil`) demonstrating packages, modules, and project layout. See its own [readme.md](36_packages_modules_project/readme.md).

### Web development

- [37_web_development/](37_web_development/) — HTTP fundamentals, servers, clients, REST APIs, JSON APIs, routing, middleware, authentication concepts, request validation, and API error handling.

### Databases

- [38_database/](38_database/) — SQL basics, `database/sql`, PostgreSQL, CRUD operations, transactions, connection pooling, and database migrations.

### Testing and production

- [39_testing_production/](39_testing_production/) — Unit testing, table-driven tests, benchmarks, test coverage, the race detector, `go vet`, logging, configuration, graceful shutdown, profiling, cross-compilation, and a sample [deployment](39_testing_production/11_deployment/) (Dockerfile + app).

## How to run

From the project directory, run a top-level lesson file with:

```bash
go run 1_hello.go
```

For files inside a numbered subdirectory (e.g. `37_web_development`, `38_database`, `39_testing_production`), `cd` into that directory first, since some of them have their own `go.mod`:

```bash
cd 37_web_development
go run 02_http_server.go
```

For the multi-package project in `36_packages_modules_project`, run the app from that directory:

```bash
cd 36_packages_modules_project
go run ./cmd/app
```

You can replace the filename/path with any lesson you want to try.

## Prerequisites

Make sure Go is installed on your machine.

- Download: https://go.dev/dl/
- Verify installation:

```bash
go version
```

## Learning path

Start from the lowest numbered file and move upward in order. Each file (or subdirectory) is designed to introduce one concept at a time, progressing from core language features to concurrency, real-world project structure, web development, databases, and production-readiness practices.

## Notes

This repository is intended for learning and experimentation. You can modify the example files to test different values and behavior.
