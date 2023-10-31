package squll

import (
	"fmt"
	"strings"
	"text/template"
)

type placeholderFormat rune

const (
	placeholderFormatDollar   placeholderFormat = '$'
	placeholderFormatQuestion placeholderFormat = '?'
	placeholderFormatColon    placeholderFormat = ':'
)

type Template struct {
	parsed      *template.Template
	placeholder placeholderFormat
}

const functionName = "argument"

func (t *Template) Build(data any) (string, []any, error) {
	var (
		args     []any
		argIndex int
	)

	tmpl, err := t.parsed.Clone()
	if err != nil {
		return "", nil, fmt.Errorf("clone: %w", err)
	}

	if t.placeholder == placeholderFormatQuestion {
		tmpl = tmpl.Funcs(template.FuncMap{
			functionName: func(arg any) string {
				args = append(args, arg)
				return string(t.placeholder)
			},
		})
	} else {
		tmpl = tmpl.Funcs(template.FuncMap{
			functionName: func(arg any) string {
				args = append(args, arg)
				argIndex++
				return fmt.Sprintf("%c%d", t.placeholder, argIndex)
			},
		})
	}

	query := &strings.Builder{}
	if err := tmpl.Execute(query, data); err != nil {
		return "", nil, fmt.Errorf("execute: %w", err)
	}
	return query.String(), args, nil
}

type config struct {
	Placeholder placeholderFormat
}

func WithColonPlaceholder() func(*config) {
	return func(c *config) {
		c.Placeholder = placeholderFormatColon
	}
}

func WithDollarPlaceholder() func(*config) {
	return func(c *config) {
		c.Placeholder = placeholderFormatDollar
	}
}

func WithQuestionPlaceholder() func(*config) {
	return func(c *config) {
		c.Placeholder = placeholderFormatQuestion
	}
}

var DefaultPlaceholderOption = WithDollarPlaceholder()

func Must(text string, options ...func(*config)) *Template {
	cfg := &config{}
	DefaultPlaceholderOption(cfg)
	for _, opt := range options {
		opt(cfg)
	}

	tmpl := template.New("")
	tmpl.Funcs(template.FuncMap{
		functionName: nop,
	})
	var err error
	tmpl, err = tmpl.Parse(text)
	if err != nil {
		panic(err)
	}

	return &Template{parsed: tmpl, placeholder: cfg.Placeholder}
}

func nop(any) string {
	return ""
}
