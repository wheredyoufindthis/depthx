# depth 

[![GoDoc](https://godoc.org/github.com/KyleBanks/depth?status.svg)](https://godoc.org/github.com/KyleBanks/depth)&nbsp; 
[![Build Status](https://travis-ci.org/KyleBanks/depth.svg?branch=master)](https://travis-ci.org/KyleBanks/depth)&nbsp;
[![Go Report Card](https://goreportcard.com/badge/github.com/KyleBanks/depth)](https://goreportcard.com/report/github.com/KyleBanks/depth)

`depth` is tool to retrieve and visualize Go source code dependency trees.

## Install

```sh
go get github.com/KyleBanks/depth/cmd/depth
```

## Usage

`depth` can be used as a standalone command-line application, or as a package within your own project.

### Command-Line

Simply execute `depth` with one or more package names to visualize:

```sh
$ depth github.com/KyleBanks/depth/cmd/depth
github.com/KyleBanks/depth/cmd/depth
  ├ encoding/json
  ├ flag
  ├ fmt
  ├ io
  ├ log
  ├ os
  ├ strings
  └ github.com/KyleBanks/depth
    ├ fmt
    ├ go/build
    ├ path
    ├ sort
    └ strings
```

You can also use `depth` on the Go standard library:

```sh
$ depth strings
strings
  ├ errors
  ├ io
  ├ unicode
  └ unicode/utf8
```

Visualizing multiple packages at a time is supported by simply naming the packages you'd like to visualize:

```sh
$ depth strings github.com/KyleBanks/depth
strings
  ├ errors
  ├ io
  ├ unicode
  └ unicode/utf8
github.com/KyleBanks/depth
  ├ fmt
  ├ go/build
  ├ path
  ├ sort
  └ strings
```

#### `-internal`

By default, `depth` only resolves the top level of dependencies for standard library packages, however you can use the `-internal` flag to visualize all internal dependencies:

```sh
$ depth -internal strings
strings
  ├ errors
  ├ io
    ├ errors
    └ sync
      ├ internal/race
        └ unsafe
      ├ runtime
        ├ runtime/internal/atomic
          └ unsafe
        ├ runtime/internal/sys
        └ unsafe
      ├ sync/atomic
        └ unsafe
      └ unsafe
  ├ unicode
  └ unicode/utf8
```

#### `-max` 

The `-max` flag limits the dependency tree to the maximum depth provided. For example, if you supply `-max 1` on the `depth` package, your output would look like so:

```
$ depth -max 1 github.com/KyleBanks/depth/cmd/depth
github.com/KyleBanks/depth/cmd/depth
  ├ encoding/json
  ├ flag
  ├ fmt
  ├ io
  ├ log
  ├ os
  ├ strings
  └ github.com/KyleBanks/depth
```

The `-max` flag is particularly useful in conjunction with the `-internal` flag which can lead to very deep dependency trees.

#### `-test`

By default, `depth` ignores dependencies that are only required for testing. However, you can view test dependencies using the `-test` flag:

```sh
$ depth -test strings
strings
  ├ bytes
  ├ errors
  ├ fmt
  ├ io
  ├ io/ioutil
  ├ math/rand
  ├ reflect
  ├ sync
  ├ testing
  ├ unicode
  ├ unicode/utf8
  └ unsafe
```

#### `-json`

The `-json` flag instructs `depth` to output dependencies in JSON format:

```sh
$ depth -json github.com/KyleBanks/depth/cmd/depth
{
  "name": "github.com/KyleBanks/depth/cmd/depth",
  "deps": [
    {
      "name": "encoding/json",
      "deps": null
    },
    ...
    {
      "name": "github.com/KyleBanks/depth",
      "deps": [
        {
          "name": "go/build",
          "deps": null
        },
        ...
      ]
    }
  ]
}
```

### Integrating With Your Project

The `depth` package can easily be used to retrieve the dependency tree for a particular package in your own project. For example, here's how you would retrieve the dependency tree for the `strings` package:

```go
import "github.com/KyleBanks/depth"

var t depth.Tree
err := t.Resolve("strings")
if err != nil {
    log.Fatal(err)
}

// Output: "'strings' has 4 dependencies."
log.Printf("'%v' has %v dependencies.", t.Root.Name, len(t.Root.Deps)) 
```

## Author

`depth` was developed by [Kyle Banks](https://twitter.com/kylewbanks).

## License

`depth` is available under [MIT](./LICENSE)