package service

import (
	"context"
	"errors"
	"testing"

	"github.com/5imonGustafsson/grpc-iot-concept/server/mocks"
	"github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

func TestService_SendWaterSoilLevel(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		var gotOrg string
		var gotBucket string
		var gotPoints []*write.Point

		s := New(&mocks.InfluxClient{
			WriteAPIBlockingFunc: func(org, bucket string) api.WriteAPIBlocking {
				gotOrg = org
				gotBucket = bucket
				return &mocks.WriteAPI{
					WritePointFunc: func(context context.Context, point ...*write.Point) error {
						gotPoints = point
						return nil
					},
				}
			},
		}, "foo-org")

		gotResponse, err := s.SendWaterSoilLevel(context.Background(), &messages.WaterSoilMetrics{
			DeviceId:      "foo-id",
			MessageId:     "bar-id",
			MoistureLevel: 1.0,
		})
		if err != nil {
			t.Errorf("unexpected error %v", err)
		}

		if got, want := gotOrg, "foo-org"; got != want {
			t.Errorf("Got organisation %s want: %s", got, want)
		}

		if got, want := gotBucket, "hydrophonic"; got != want {
			t.Errorf("Got bucket %s want: %s", got, want)
		}

		if got, want := len(gotPoints), 1; got != want {
			t.Errorf("Got batch point length %v want %v", got, want)
		}

		wantResponse := messages.MetricsReply{
			MessageId:  "bar-id",
			StatusCode: 200,
		}

		if got, want := gotResponse, &wantResponse; got.MessageId != want.MessageId || got.StatusCode != want.StatusCode {
			t.Errorf("got response: %v, want: %v", gotResponse, &wantResponse)
		}

	})

	t.Run("Error", func(t *testing.T) {

		fooErr := errors.New("foo error")

		s := New(&mocks.InfluxClient{
			WriteAPIBlockingFunc: func(org, bucket string) api.WriteAPIBlocking {
				return &mocks.WriteAPI{
					WritePointFunc: func(context context.Context, point ...*write.Point) error {
						return fooErr
					},
				}
			},
		}, "foo-org")

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
