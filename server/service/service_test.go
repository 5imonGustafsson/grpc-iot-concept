package service

import (
	"context"
	"errors"
	"testing"

	"github.com/5imonGustafsson/grpc-iot-concept/server/mocks"
	"github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	influx "github.com/influxdata/influxdb/client/v2"
)

func TestService_SendWaterSoilLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotBatchPoints []*influx.Point

		s := New(&mocks.InfluxClient{
			WriteCtxFunc: func(ctx context.Context, bp influx.BatchPoints) error {
				gotBatchPoints = bp.Points()
				return nil
			},
		})

		wantResponse := messages.MetricsReply{
			MessageId:  "bar-id",
			StatusCode: 200,
		}

		gotResponse, err := s.SendWaterSoilLevel(context.Background(), &messages.WaterSoilMetrics{
			DeviceId:      "foo-id",
			MessageId:     "bar-id",
			MoistureLevel: 1.0,
		})
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if got, want := len(gotBatchPoints), 1; got != want {
			t.Errorf("Expected batch point length %v got %v", want, got)
		}

		if got, want := gotResponse, &wantResponse; got.MessageId != want.MessageId || got.StatusCode != want.StatusCode {
			t.Errorf("expected response: %v, got: %v", &wantResponse, gotResponse)
		}

	})

	t.Run("Error", func(t *testing.T) {

		fooErr := errors.New("foo error")

		s := New(&mocks.InfluxClient{
			WriteCtxFunc: func(ctx context.Context, bp influx.BatchPoints) error {
				return fooErr
			},
		})

		_, err := s.SendWaterSoilLevel(context.Background(), &messages.WaterSoilMetrics{
			DeviceId:      "foo-id",
			MessageId:     "bar-id",
			MoistureLevel: 1.0,
		})
		if err == nil {
			t.Errorf("error is nil")
		}

		if got, want := err, fooErr; !errors.Is(got, want) {
			t.Errorf("unexpected error: %v", got)
		}

	})
}
