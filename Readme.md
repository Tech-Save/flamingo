# Flamingo

## What is Flamingo

Flamingo is a go based, opinionated framework for pluggable web projects.

## Setup

In order to work properly, Flamingo needs to be checked out at the following location: `$GOPATH/src/go.aoe.com/flamingo`
```sh
git clone git@gitlab.aoe.com:shared/flamingo/flamingo.git $GOPATH/src/go.aoe.com/flamingo
```

## Necessary tooling

### dep

Dependency manager:

`go get -u github.com/golang/dep/cmd/dep`

Usage

`dep ensure` or `dep ensure -vendor-only`

### go-bindata

Static file compiler for fakeservices, etc

`go get -u github.com/jteeuwen/go-bindata/...`

## Docs

To read the documentation:

```
make docs
```

This will start building and previewing the mkdocs based documentation in a Docker container.

To view the docs open  [Docs](http://localhost:8000)

-----------------

# Mockery (to create Mocks)

https://github.com/vektra/mockery

Usage:
mockery -name <Name of Interface>



What we changed

We forked https://github.com/leebenson/conform into https://github.com/AOEpeople/conform and removed the faulties regarding stripping special chars in a wrong way.


special chars are simply removed from the given string
apostrophes (') are allowed in a name e.g. Andy O'Dowell is possible now
dashes and spaces are shrinked together e.g. " ab   -  cd" becomes "Ab-Cd"
