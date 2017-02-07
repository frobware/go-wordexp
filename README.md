# go-wordexp

This is a Cgo wrapper around wordexp(3) and explicitly passes
WRDE_NOCMD | WRDE_UNDEF ("don't do command substitution", "treat
undefined variables as errors") as flags to wordexp(3).

## CI status

 * Travis: [![Build Status](https://travis-ci.org/frobware/go-wordexp.svg?branch=master)](https://travis-ci.org/frobware/go-wordexp)
 * Code Coverage: [![Coverage Status](http://codecov.io/github/frobware/go-wordexp/coverage.svg?branch=master)](http://codecov.io/github/frobware/go-wordexp?branch=master)
 
