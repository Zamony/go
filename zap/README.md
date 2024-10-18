Minimal type-safe wrapper around the standard `log/slog` library:
* no external dependencies;
* provides error attribute with stacktrace.

```
type Logger struct {
	Logger *slog.Logger
}
func (l *Logger) Debug(msg string, attrs ...slog.Attr)
func (l *Logger) Info(msg string, attrs ...slog.Attr)
func (l *Logger) Warn(msg string, attrs ...slog.Attr)
func (l *Logger) Error(msg string, attrs ...slog.Attr)
func (l *Logger) With(attrs ...slog.Attr) Logger

var ErrorStackMarshaler func(error) string
func Error(err error) slog.Attr 
```
