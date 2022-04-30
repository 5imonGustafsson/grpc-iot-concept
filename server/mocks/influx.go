package mocks

import (
	"context"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type InfluxClient struct {
	WriteAPIBlockingFunc func(org string, bucket string) api.WriteAPIBlocking
	influxdb2.Client
}

type WriteAPI struct {
	WritePointFunc func(ctx context.Context, point ...*write.Point) error
	api.WriteAPIBlocking
}

func (c *InfluxClient) WriteAPIBlocking(org string, bucket string) api.WriteAPIBlocking {
	return c.WriteAPIBlockingFunc(org, bucket)
}

func (w *WriteAPI) WritePoint(ctx context.Context, point ...*write.Point) error {
	return w.WritePointFunc(ctx, point...)
}
