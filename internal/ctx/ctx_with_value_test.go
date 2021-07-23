package ctx

import (
	"context"
	"fmt"
	"testing"
)

func Test_ctx_with_value(t *testing.T) {
	ProcessRequest("da","opa")
}


func ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), "userID", userID)
	ctx = context.WithValue(ctx, "authToken", authToken)
	HandleResponse(ctx)
}

func HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (%v)",
		ctx.Value("userID"),
		ctx.Value("authToken"),
	)
}

func TestSafeUseCtxValue(t *testing.T) {
	
}

type ctxKey int

const (
	ctxUserID ctxKey = iota
	ctxAuthToken
)

func UserID(c context.Context) string {
	return c.Value(ctxUserID).(string)
}

func AuthToken(c context.Context) string {
	return c.Value(ctxAuthToken).(string)
}

func _ProcessRequest(userID, authToken string) {
	ctx := context.WithValue(context.Background(), ctxUserID, userID)
	ctx = context.WithValue(ctx, ctxAuthToken, authToken)
	HandleResponse(ctx)
}

func _HandleResponse(ctx context.Context) {
	fmt.Printf(
		"handling response for %v (auth: %v)",
		UserID(ctx),
		AuthToken(ctx),
	)
}