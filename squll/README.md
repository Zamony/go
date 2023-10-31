Simple SQL templating:
* configurable placeholders (`$1`, `:2`, `?`);
* no external dependencies;

**API**
```go

import "github.com/Zamony/go/squll"

var queryInsertScores = squll.Must(`
    INSERT INTO scores VALUES
    {{range $i, $v := .}}
        {{if $i}}, {{end}}
        ({{argument $v.Name}}, {{argument $v.Score}})
    {{end}}
`)

func insertScores() error {
    query, args, err := queryInsertScores.Build([]UserScores{
        {"Alice", 10},
        {"Bob",   20},
        {"Chris", 30},
    })
    if err != nil {
        return err
    }
    // Query:
    // INSERT INTO scores VALUES ($1, $2), ($3, $4), ($5, $6)

    // Args:
    // []any{"Alice", 10, "Bob", 20, "Chris", 30}
}
```
