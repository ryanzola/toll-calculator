package client

import (
	"context"

	"github.com/ryanzola/toll-calculator/types"
)

type Client interface {
	Aggregate(context.Context, *types.AggregateRequest) error
}
