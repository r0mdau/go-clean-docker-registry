go-clean-docker-registry
========================

Super simple cli tool written in Go to clean your v2 docker registry.

File `src/main.go` is the entrypoint of the program.

## Prerequisites
- go
- make

## Quick start

Multiple flags are settable, mandatory are image name and registry api url.
All flags have a shortcut alias.

Two main commands :
- show : show all tags associated with image
- delete : delete all tags according to provided flags to the program

### Use it

To use this project you can simply use `go run` or launch the binary.

Show all images in registry :

    go run main.go cmd.go registry.go showimages -u https://registry.docker.example.com

Show all tags of specified image :

    go run main.go cmd.go registry.go showtags -u https://registry.docker.example.com -i r0mdau/nodejs

### Build
Command `make` to build amd64 binary.
```
make
```

### (Un)Install

```
make install
make uninstall
```

## TODO
- [ ] Makefile
- [ ] Readme
- [ ] Load flags using a yaml config file
- [ ] Logic for cli flags
- [ ] get tags from registry
- [ ] if dryrun output tags to be deleted
- [ ] else delete docker images:tag
- [ ] docker hub api authent (JWT) : https://hub.docker.com/support/doc/how-do-i-authenticate-with-the-v2-api
