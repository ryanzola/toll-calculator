package main

import (
	"context"
	"log"
	"time"

	"github.com/ryanzola/toll-calculator/aggregator/client"
	"github.com/ryanzola/toll-calculator/types"
)

func main() {
	c, err := client.NewGRPCClient(":3003")
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Aggregate(context.Background(), &types.AggregateRequest{
		ObuID: 1,
		Value: 58.55,
		Unix:  time.Now().UnixNano(),
	}); err != nil {
		log.Fatal(err)
	}
}
