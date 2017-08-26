# Cache

## A minimalistic Time aware Least Recently Used (TLRU) cache implementation.

If you need something quick and simple to use in your tests, or you don't want to install another component just to keep some hot data in the memory, this package is for you.

It is not intended to be used as a storage. To keep your application from eating up your memory, you can set limit on number of items, or number of stored bytes.

Performance was not a major driving force behind this package, but there is not a lot going on so you can expect fair performance. Feel free to benchmark it and make a PR. ;)

__A word about Time awareness:__ entries in the cache are not actively deleted based on TTL. The TTL of an entry is evaluated upon accessing the entry; expired entries are deleted only on an attempt to retrieve them. In other words, if nothing fetches an entry, it will only be deleted when it becomes the least recently used entry, which is inevitable eventually. As a result of this, the cache is only time aware in terms of keeping stale data inaccessible for fetch requests, but it's not actively freeing up memory.

## Installation

`go get github.com/zgiber/cache`

## Usage

See the examples directory.


