go-clean-docker-registry
========================

[![Go Report Card](https://goreportcard.com/badge/github.com/r0mdau/go-clean-docker-registry)](https://goreportcard.com/report/github.com/r0mdau/go-clean-docker-registry)

Super simple cli tool written in Go to clean your v2 docker registry.

File `main.go` is the entrypoint of the program.

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

Show all images in the registry :

    go run main.go showimages -u https://registry.docker.example.com
    # or
    go-clean-docker-registry showimages -u https://registry.docker.example.com

Show all tags of specified image :

    go-clean-docker-registry showtags -u https://registry.docker.example.com -i r0mdau/nodejs

Delete all tags of specified image :

    go-clean-docker-registry delete -u https://registry.docker.example.com -i r0mdau/nodejs

Delete all matched tags of specified image and keep the 10 last, semver versioning for sorting :

    go-clean-docker-registry delete -u https://registry.docker.example.com -i r0mdau/nodejs -t master-* -k 10

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
- [ ] Load flags using a yaml config file
- [ ] docker hub api authent (JWT) : https://hub.docker.com/support/doc/how-do-i-authenticate-with-the-v2-api
- [ ] be satisfied with code quality and code coverage
- [x] Makefile
- [x] Logic for cli flags
- [x] get images from registry
- [x] get tags from registry
- [x] DELETE action, if dryrun output tags to be deleted
- [x] DELETE action, else delete docker images:tag
- [x] folder structure with go packages
