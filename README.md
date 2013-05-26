Rotten tomato API for go (golang)
=====================================

Uses the [API](http://developer.rottentomatoes.com/docs) to access
movie data from rotten tomatoes.


Setup
-----

The API key is missing as that information is private. Create a file (i.e. key.go)
that contanis your (private) api key as follows:

```go
package rotten_tomatoes

const rtkey := "<your api key>"
```
