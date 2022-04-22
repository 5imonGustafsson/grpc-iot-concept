package mocks

import (
	"context"

	influx "github.com/influxdata/influxdb/client/v2"
)

type InfluxClient struct {
	WriteCtxFunc func(ctx context.Context, bp influx.BatchPoints) error
	influx.Client
}

func (c *InfluxClient) WriteCtx(ctx context.Context, bp influx.BatchPoints) error {
	return c.WriteCtxFunc(ctx, bp)
}
