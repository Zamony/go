Сonvert http code to its type

**Example**

```go
package main

import (
    "fmt"
    "net/http"

    "github.com/Zamony/go/httpstatus"
)

func main() {
    httpstatus.From(http.StatusOK) // httpstatus.Success
}

```
