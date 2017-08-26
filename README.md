# Cache

## A minimalistic Time aware Least Recently Used (TLRU) cache implementation.

If you need something quick and simple to use in your tests, or you don't want to install another component just to keep some hot data in the memory, this package is for you.

It is not intended to be used as a storage. To keep your application from eating up your memory, you can set limit on number of items, or number of stored bytes.

Performance was not a major driving force behind this package, but there is not a lot going on so you can expect fair performance. Feel free to benchmark it and make a PR. ;)

## Installation

`go get github.com/zgiber/cache`

## Usage

See the examples directory.


