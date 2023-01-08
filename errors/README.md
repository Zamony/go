Simple error handling primitives.

**Features**
* stack traces;
* errors wrapping;
* compatible with stdlib;
* no external dependencies;

**Example**

```go
package main

import (
	"github.com/Zamony/go/errors"
	"github.com/Zamony/go/errors/stackerr"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {
	zerolog.ErrorStackMarshaler = stackerr.MarshalCompactZerolog
	err := errors.New("nasty error")
	log.Error().Stack().Err(err).Msg("Main has failed")
}

// Output: {"level":"error","stack":"main.main:12/runtime.main/goexit","error":"nasty error","time":"2023-01-08T14:21:15+03:00","message":"Main has failed"}

```
