# GEE (Go Error Expander)

Go is stupid simple. Unfortunately, some parts of it are plain stupid. One of them is error handling.

The `gee` CLI tool simply expands placeholders into error handling, saving the programmer from typing:

```go
...
if err != nil {
    return err
}
...
```
over and over again, countless times across the codebase.

## Getting started

- Add `import fmt` in the file if errors need to be wrapped with custom messages
- Add a `var _ error` declaration at the top of the function
- Use `_` in place of an error
- Add a `//gee:` comment above the line with the error. If a message is present after the colon, GEE wraps the error with that message and returns it. If not, it directly returns the error.

GEE will automatically replace the placeholder `_`s and inject error handling.

Consider a sample `input.go` file:

```go
package main

import "fmt"

func foo() error {
    return nil
}

func main() {
    var _ error
    //gee:
    _ = foo()
    //gee:failed to run foo
    _ = foo()
}
```

Upon running `gee input.go`, GEE will expand this into:

```go
package main

import "fmt"

func foo() error {
    return nil
}

func main() {
    var err error
    err = foo()
    if err != nil {
        return err
    }
    err = foo()
    if err != nil {
        return fmt.Errorf("failed to run foo: %w", err)
    }
}
```

## Installation

The easiest way is via `go install`:

```bash
go install github.com/snprajwal/gee
```

## Usage

```bash
# Providing a directory as the argument runs the CLI on all `.go` files in the directory
gee cli/
# Providing a single `.go` file runs the CLI on just that file
gee main.go
```
By default, GEE prints the expanded Go code to `stdout`. To modify the original `.go` files, run the CLI with the `-i` or `--in-place` flag.

```bash
# Providing the `-i` or `--in-place` flag modifies the file in-place
gee -i main.go
```

## License

This project is licensed under the MIT License and is available [here](LICENSE).
