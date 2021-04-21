go-clean-docker-registry
========================

Super simple cli tool written in Go to clean your v2 docker registry.

File `src/main.go` is the entrypoint of the program.

## Prerequisites
- go
- make

## Quick start

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
- [ ] Logic for cli flags
- [ ] get tags from registry
- [ ] if dryrun output tags to be deleted
- [ ] else delete docker images:tag
