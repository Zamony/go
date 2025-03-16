package sqlbatch_test

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/Zamony/go/sqlbatch"
)

func TestBatch(t *testing.T) {
	batch := sqlbatch.New(nil)
	defer batch.Close()

	persons := []Person{
		{"Ivan", 33},
		{"Alexey", 45},
		{"Maria", 25},
	}
	for _, person := range persons {
		batch.Append(person.Name, person.Age)
	}

	gotQuery := batch.BuildQuery("INSERT INTO example VALUES ", "")
	gotArgs := batch.BuildArguments()
	wantQuery := "INSERT INTO example VALUES ($1,$2),($3,$4),($5,$6);"
	if wantQuery != gotQuery {
		t.Errorf("\n+%s\n-%s\n", gotQuery, wantQuery)
	}
	wantArgs := []any{"Ivan", 33, "Alexey", 45, "Maria", 25}
	if !reflect.DeepEqual(gotArgs, wantArgs) {
		t.Errorf("\n+%v\n-%v\n", gotArgs, wantArgs)
	}
}

func TestBatchWithDifferentPlaceholderFormats(t *testing.T) {
	tests := []struct {
		name     string
		format   byte
		expected string
	}{
		{"Dollar", sqlbatch.PlaceholderFormatDollar, "INSERT INTO example VALUES ($1,$2),($3,$4);"},
		{"Colon", sqlbatch.PlaceholderFormatColon, "INSERT INTO example VALUES (:1,:2),(:3,:4);"},
		{"Question", sqlbatch.PlaceholderFormatQuestion, "INSERT INTO example VALUES (?,?),(?,?);"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			batch := sqlbatch.New(&sqlbatch.Options{
				PlaceholderFormat: tt.format,
			})
			defer batch.Close()

			batch.Append("Ivan", 33)
			batch.Append("Alexey", 45)

			gotQuery := batch.BuildQuery("INSERT INTO example VALUES ", "")
			if gotQuery != tt.expected {
				t.Errorf("\n+%s\n-%s\n", gotQuery, tt.expected)
			}
		})
	}
}

func TestBatchWithSuffix(t *testing.T) {
	batch := sqlbatch.New(nil)
	defer batch.Close()

	batch.Append("Ivan", 33)
	batch.Append("Alexey", 45)

	gotQuery := batch.BuildQuery("INSERT INTO example VALUES ", "ON CONFLICT DO NOTHING")
	wantQuery := "INSERT INTO example VALUES ($1,$2),($3,$4) ON CONFLICT DO NOTHING;"
	if wantQuery != gotQuery {
		t.Errorf("Expected query with suffix to be %q, got %q", wantQuery, gotQuery)
	}
}

func TestBatchConcurrency(*testing.T) {
	var dummy atomic.Int64 // avoid compiler optimization
	var wg sync.WaitGroup
	const count = 1000
	persons := make([]Person, count)
	for i := 1; i < count; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			batch := sqlbatch.New(nil) // same cache for all queries
			defer batch.Close()

			for _, person := range persons[:i] {
				batch.Append(person.Name, person.Age)
			}

			query := batch.BuildQuery("INSERT INTO example VALUES ", "")
			dummy.Add(int64(len(query)))
		}()
	}

	wg.Wait()
}

func BenchmarkBuildQueryBatch(b *testing.B) {
	persons := makePersons(10000)
	for b.Loop() {
		batch := sqlbatch.New(nil)
		for _, person := range persons {
			batch.Append(person.Name, person.Age)
		}
		_ = batch.BuildQuery("INSERT INTO example VALUES ", "")
		batch.Close()
	}
}

func BenchmarkBuildQueryBatchCached(b *testing.B) {
	persons := makePersons(10000)
	cache := make(map[int]string)
	for b.Loop() {
		batch := sqlbatch.New(nil)
		for _, person := range persons {
			batch.Append(person.Name, person.Age)
		}
		if _, ok := cache[len(persons)]; !ok {
			query := batch.BuildQuery("INSERT INTO example VALUES ", "")
			cache[len(persons)] = query
		}
		batch.Close()
	}
}

func BenchmarkBuildQueryStringsBuilder(b *testing.B) {
	persons := makePersons(10000)
	for b.Loop() {
		var builder strings.Builder
		builder.WriteString("INSERT INTO example VALUES ")
		args := make([]any, 0, len(persons)*2)
		count := 1
		for _, person := range persons {
			fmt.Fprintf(&builder, "($%d,$%d),", count, count+1)
			count += 2
			args = append(args, person.Name, person.Age)
		}
		_, _ = builder.String(), args
	}
}

type Person struct {
	Name string
	Age  int
}

func makePersons(size int) []Person {
	persons := make([]Person, size)
	for i := range persons {
		persons[i].Name = "Person" + strconv.Itoa(i)
		persons[i].Age = i
	}
	return persons
}
