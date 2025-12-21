package atylabmongo

import (
	"context"
	"time"
)

type MongoCtxSvc struct {
	Ctx    context.Context
	Cancel context.CancelFunc
}

func NewMongoCtxSvc() *MongoCtxSvc {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	return &MongoCtxSvc{
		Ctx:    ctx,
		Cancel: cancel,
	}
}

func NewMongoCtxSvcWithTimeout(timeout time.Duration) *MongoCtxSvc {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	return &MongoCtxSvc{
		Ctx:    ctx,
		Cancel: cancel,
	}
}
