package squll_test

import (
	"reflect"
	"testing"

	"github.com/Zamony/go/squll"
)

func TestBuildQuery(t *testing.T) {
	template := squll.Must(`INSERT INTO scores VALUES ({{argument .Name}}, {{argument .Score}})`)

	type UserScore struct {
		Name  string
		Score int
	}

	params := UserScore{"Alice", 10}
	wantQuery := `INSERT INTO scores VALUES ($1, $2)`
	wantArgs := []any{"Alice", 10}

	t.Run("fresh template", func(t *testing.T) {
		checkQuery(t, template, params, wantQuery, wantArgs)
	})
	t.Run("used template", func(t *testing.T) {
		checkQuery(t, template, params, wantQuery, wantArgs)
	})
}

func TestBuildQueryStaticPlaceholder(t *testing.T) {
	template := squll.Must(
		`INSERT INTO names VALUES ({{argument .}})`,
		squll.WithQuestionPlaceholder(),
	)

	name := "Alice"
	wantQuery := `INSERT INTO names VALUES (?)`
	checkQuery(t, template, name, wantQuery, []any{name})
}

func checkQuery(t *testing.T, template *squll.Template, params any, wantQuery string, wantArgs []any) {
	query, args, err := template.Build(params)
	if err != nil {
		t.Fatal(err)
	}

	if query != wantQuery {
		t.Fatalf("Query mismatch\nWant %s\nGot %s", wantQuery, query)
	}

	if !reflect.DeepEqual(args, wantArgs) {
		t.Fatalf("Args mismatch\nWant %v\nGot %v", wantArgs, args)
	}
}
