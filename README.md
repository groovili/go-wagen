# go-wagen
[![Go Report Card](https://goreportcard.com/badge/github.com/groovili/go-wagen)](https://goreportcard.com/report/github.com/groovili/go-wagen)

Single binary **web application generator** for Go. 

**Creates project boilerplate** and gives the ability to focus on the implementation of app functionality.

#### Generates:
 * project layout, according to community best practices
 * Makefile for wrapping project related routines 
 * Dockerfile and docker-compose for local development
 * containers for code, tests and [golangci-lint](https://github.com/golangci/golangci-lint)
 * configuration management with [viper](https://github.com/spf13/viper)
 * logging with [logrus](http://github.com/sirupsen/logrus) or [zap](https://github.com/uber-go/zap)
 * routing with [gorilla/mux](https://github.com/gorilla/mux) or [chi](github.com/go-chi/chi)
 * default health check, http handler and logger middleware
 
 **go-wagen is a starter pack** for typical web application. 
 
 It doesn't aim to generalize the whole project workflow, push framework, or architecture.
 Consists of framework-agnostic components that are common in most of the applications.
 
 All contributions, issues, requests or feedback are warmly welcome.
 
 ## Installation
 
 Install pre-built binary on [releases page](https://github.com/groovili/go-wagen/releases):
 
 1. `curl -LJO https://github.com/groovili/go-wagen/releases/download/latest/go-wagen-osx.tar.gz`
 2. `tar -f go-wagen-osx.tar.gz -x`
 2. `./go-wagen --path=/absoule/path/to/project`
 
 or build from source code:
 
 1. `git clone https://github.com/groovili/go-wagen && cd go-wagen`
 2. `make install` - will install [packr](https://github.com/gobuffalo/packr) to wrap templates to binary
 3. `make build`
 4. `./go-wagen --path=/absoule/path/to/project`
 
 Binary doesn't need to be in `$GOPATH` and works without any dependencies.
 
 ## Usage
 
 Video instruction soon will be here.
 
 1. `./go-wagen --path=/absoule/path/to/project` and select dependencies
 2. `cd /absoule/path/to/project`
 3. `go mod vendor`
 4. `make run` - will build and run container with code
 5. `make test` - to run container with tests
 6. `make lint` - to run linter for source code
 
 #### Structure
 
 Application generated with **go-wagen** will have following structure:
 ```
 ├── Makefile
 ├── cmd
 │   └── example
 │       └── example.go - app entrypoint
 ├── config - config presets for each env of default pipline
 │   ├── app.dev.yml
 │   ├── app.local.yml
 │   └── app.yml
 ├── deploy
 │   ├── Dockerfile
 │   └── docker-compose.yml
 ├── go.mod
 ├── internal - for shared internal tools
 ├── server
 │   ├── handlers
 │   │   ├── hello.go
 │   │   └── ping.go
 │   └── middleware
 │       └── log.go
 ├── storage - db/storages
 └── vendor - downloaded modules
```