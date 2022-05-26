package userid

import "context"

// userIdKey is used to store the user ID in the request context.
// Authentication middleware can set this key.
var userIdKey = &contextKey{"userID"}

type contextKey struct {
	name string
}

// Get the user ID from the context.
func Get(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	userID, ok := ctx.Value(userIdKey).(*string)
	if !ok {
		return ""
	}
	if userID == nil {
		return ""
	}
	return *userID
}

// Set the user ID in the context.
// If userid.Init() has been called previously,
// this sets the existing user ID in the context.
func Set(ctx context.Context, uid string) context.Context {
	// if there is a user ID already in the context, set the existing one.
	if userID, ok := ctx.Value(userIdKey).(*string); ok {
		*userID = uid
		return ctx
	}

	return context.WithValue(ctx, userIdKey, &uid)
}

// Init sets up an empty user ID in the context.
func Init(ctx context.Context) context.Context {
	uid := ""
	return context.WithValue(ctx, userIdKey, &uid)
}
