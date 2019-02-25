<h1 align="center">Mini Web Crawler</h1>

<p align="center">
    <a href="https://codecov.io/gh/TimTosi/mcrawler">
        <img src="https://codecov.io/gh/TimTosi/mcrawler/branch/master/graph/badge.svg" alt="codecov" />
    </a>
    <a href="https://circleci.com/gh/TimTosi/mcrawler">
        <img src="https://circleci.com/gh/TimTosi/mcrawler.svg?style=shield" alt="CircleCI" />
    </a>
    <a href="https://goreportcard.com/report/github.com/timtosi/mcrawler">
        <img src="https://goreportcard.com/badge/github.com/timtosi/mcrawler" alt="Go Report" />
    </a>
    <a href="https://godoc.org/github.com/timtosi/mcrawler">
        <img src="https://godoc.org/github.com/timtosi/mcrawler?status.svg" alt="GoDoc" />
    </a>
    <a href="https://opensource.org/licenses/MIT">
        <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" />
    </a>
</p>

## Table of Contents
- [About](#what-this-repository-is-about)
- [Architecture](#architecture)
- [Quickstart with Docker](#quickstart-with-docker)
- [Quickstart with binary](#quickstart-without-docker)
- [Component List](#component-list)
- [How To Add a Component](#how-to-add-a-component)
- [FAQ](#faq)
- [Support & Feedbacks](#not-good-enough-)


## What this repository is about

This repository contains source code of an implementation of a simple web
crawler written in Go.

This work has been inspired from the [Mercator](http://dl.acm.org/citation.cfm?id=598733)
web crawler, minus the distributed part.


## Architecture

This web crawler is composed of a configurable pipeline based on the
`internal.Pipe` type. Any `internal.Pipe` part of the pipeline can be removed,
their order changed or even combined with new user based `internal.Pipe`
introducing new features or filters to modify the internal behaviour of this
program.

This project comes with already provided `internal.Pipe` that you can use.

The default pipeline provided with this project and found in the
[cmd](https://github.com/TimTosi/mcrawler/blob/master/cmd main.go#L26-33)
directory looks like this:


```

S        +-------------------+       +-------------------+
T        |                   |       |                   |
A  ----> |     Archiver      |------>|      Mapper       |--+
R        |                   |       |                   |  |
T        +-------------------+       +-------------------+  |
                                                            |
                                                            |
 +----------------------------------------------------------+
 |
 |
 |      +-------------------+       +-------------------+
 |      |                   |       |                   |
 +----> |     Follower      |------>|      Worker       |---+
        |                   |       |                   |   |
        +-------------------+       +-------------------+   |
                                                            |
                                                            |
 +----------------------------------------------------------+
 |
 |
 |      +-------------------+        S
 |      |                   |        T
 +----> |     Extractor     |------> A
        |                   |        R
        +-------------------+        T

```

This pipeline cycles on itself so you have to introduce edge condition
mechanisms in order to avoid infinite loops if you do not use those provided by
[`default`](https://github.com/TimTosi/mcrawler/blob/master/internal/archiver.go).

All the goroutines are controlled & coordinated through a `sync.WaitGroup`
created in [crawler.Run](https://github.com/TimTosi/mcrawler/blob/master/internal/crawler/crawler.go#L42-58).

> :bulb: If you introduce new `internal.Pipe`s in the pipeline, don't forget
> to `wg.Done()` each time you discard an element or to `wg.Add(1)` each time
> you add a new element in the pipeline.


## Quickstart

First, go get this repository:
```sh
go get -d github.com/timtosi/mcrawler
```


### Quickstart with Docker

> :exclamation: If you don't have [Docker](https://docs.docker.com/install/) and
> [Docker Compose](https://docs.docker.com/compose/) installed, you still can
> execute this program by [compiling the binary](#quickstart-without-docker). 

This program comes with an already configured [Docker Compose](https://github.com/TimTosi/mcrawler/blob/master/deployments/docker-compose.yaml)
that crawls a website located at `localhost:8080`.

You can use the `run` target in the provided Makefile to use it easily.

> :bulb: If you want to change the crawled target, you will have to update the
> [Docker Compose file](https://github.com/TimTosi/mcrawler/blob/master/deployments/docker-compose.yaml#L10)
> accordingly.


### Quickstart without Docker

First install dependencies & compile the binary:
```sh
cd $GOPATH/src/github.com/timtosi/mcrawler/
make install && make build
```

Then launch the program by specifying the target in argument:
```sh
cd $GOPATH/src/github.com/timtosi/mcrawler/cmd/
go build && ./mcrawler "http://localhost:8080"
```

## Component List

Here is a list and description of components provided with this program:


## How To Add a Component

In order to add a component in the pipeline, you need to create a `struct`
implementing the [`internal.Pipe`](https://github.com/TimTosi/mcrawler/blob/master/internal/pipe.go#L11-13)
interface.

```golang
package example

// UserPipe is a `struct` implementaing the `internal.Pipe` interface.
type UserPipe struct {
	// properties ...
}

// NewUserPipe returns a new `example.UserPipe`.
func NewUserPipe() *UserPipe {
    return &UserPipe{}
}

// Pipe is a user defined function used in the pipeline launched by
// `crawler.Crawler`.
func (up *UserPipe) Pipe(wg *sync.WaitGroup, in <-chan *domain.Target, out chan<- *domain.Target) {
	defer close(out)

	for t := range in {
            //
            // --------> Here, do something with element received from `in`.
            //
			wg.Done() // Don't forget to wg.Done() when you discard an element !
		} else {
			out <- t
		}
	}
}
```


Then you just have to plug it in the main:

```golang
package main

import (
	"log"
	"os"

	"github.com/user/example"
	"github.com/timtosi/mcrawler/internal/crawler"
	"github.com/timtosi/mcrawler/internal/domain"
)

func main() {
	if len(os.Args[1]) == 0 {
		log.Fatal(`usage: ./mcrawler <BASE_URL>`)
	}

	t := domain.NewTarget(os.Args[1])

	if err := crawler.NewCrawler().Run(
		example.NewUserPipe(), // ------- > Insert here !!!
	); err != nil {
		log.Fatal(err)
	}
	log.Printf("shutdown")
}
```

## FAQ

None so far :raised_hands:


## License

Every file provided here is available under the [MIT License](http://opensource.org/licenses/MIT).


## Not Good Enough ?

If you encouter any issue by using what is provided here, please
[let me know](https://github.com/TimTosi/mcrawler/issues) ! 
Help me to improve by sending your thoughts to timothee.tosi@gmail.com !
