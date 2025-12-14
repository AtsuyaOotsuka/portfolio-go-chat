package atylabmongo

import (
	"context"
	"fmt"
	"time"
)

type MongoCtxSvc struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewMongoCtxSvc() *MongoCtxSvc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	fmt.Println("Created new MongoDB context")
	return &MongoCtxSvc{
		Ctx:    ctx,
		Cancel: cancel,
	}
}
