package logging

import "context"

type ctxLogger struct{}

// ContextWithLogger добавляет логгер к конекту
func ContextWithLogger(ctx context.Context, l Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// LoggerFromContext возвращает логгер из контекста
func loggerFromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*logger); ok {
		return l
	}
	return NewLogger()
}
