go-clean-docker-registry
========================

[![Go Report Card](https://goreportcard.com/badge/github.com/r0mdau/go-clean-docker-registry)](https://goreportcard.com/report/github.com/r0mdau/go-clean-docker-registry)

Super simple cli tool written in Go to clean your v2 docker registry.

(!) NOW only showimages and showtags are developed. Delete is in progress. See todolist below.

File `src/main.go` is the entrypoint of the program.

## Prerequisites
- go
- make

## Quick start

Multiple flags are settable, some are mandatory, help menu will help you.
All flags have a shortcut alias.

Three commands :
- showimages : show all images in the registry
- showtags : show all tags associated with image
- delete : delete all tags according to provided flags to the program

### Use it

To use this project you can simply use `go run` or launch the binary.

Show all images in registry :

    go run main.go showimages -u https://registry.docker.example.com
    # or
    ./go-clean-docker-registry showimages -u https://registry.docker.example.com

Show all tags of specified image :

    go run main.go showtags -u https://registry.docker.example.com -i r0mdau/nodejs

### Build
Command `make` to build amd64 binary.
```
make
# build with docker
make build-docker
```

### (Un)Install

```
make install
make uninstall
```

## TODO
- [x] Makefile
- [ ] Load flags using a yaml config file
- [x] Logic for cli flags
- [x] get images from registry
- [x] get tags from registry
- [ ] DELETE action, if dryrun output tags to be deleted
- [ ] DELETE action, else delete docker images:tag
- [ ] docker hub api authent (JWT) : https://hub.docker.com/support/doc/how-do-i-authenticate-with-the-v2-api
- [x] folder structure with go packages
