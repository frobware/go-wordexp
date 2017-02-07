[![Travis CI](https://travis-ci.org/frobware/go-wordexp.svg?branch=master)](https://travis-ci.org/frobware/go-wordexp)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/frobware/go-wordexp)
[![Coverage Status](http://codecov.io/github/frobware/go-wordexp/coverage.svg?branch=master)](http://codecov.io/github/frobware/go-wordexp?branch=master)

# go-wordexp

This is a Cgo wrapper around wordexp(3) and explicitly passes
WRDE_NOCMD | WRDE_UNDEF ("don't do command substitution", "treat
undefined variables as errors") as flags to wordexp(3).

Head to the [Go documentation](https://godoc.org/github.com/frobware/go-wordexp) to see available methods and examples.
