package glog

import "context"

type contextKey string

const (
	contextKeyData   = contextKey("data")
	contextKeyPrefix = contextKey("prefix")
)

// ContextWithData creates a context as extension of the parent context, with the provided data stored as a value.
// Any existing data on the parent context will be replaced.
func ContextWithData(ctx context.Context, args ...interface{}) context.Context {
	return context.WithValue(ctx, contextKeyData, args)
}

// ContextWithPrefix creates a context as extension of the parent context, with the provided prefix stored as a value.
// Any existing prefix on the parent context will be replaced.
func ContextWithPrefix(ctx context.Context, prefix string) context.Context {
	return context.WithValue(ctx, contextKeyPrefix, prefix)
}
