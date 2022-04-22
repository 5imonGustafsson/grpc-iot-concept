package service

import (
	"context"
	"log"
	"time"

	pb "github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	influx "github.com/influxdata/influxdb/client/v2"
)

// server is used to implement messages.IoT
type service struct {
	influxClient influx.Client
	pb.UnimplementedIoTServer
}

func New(influxClient influx.Client) pb.IoTServer {
	return &service{
		influxClient: influxClient,
	}
}

func (s *service) SendWaterSoilLevel(ctx context.Context, metrics *pb.WaterSoilMetrics) (*pb.MetricsReply, error) {
	log.Printf("Received message: %v from device: %v", metrics.GetMessageId(), metrics.GetDeviceId())

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: "hydrophonic",
	})

	if err != nil {
		return &pb.MetricsReply{
			MessageId:  metrics.GetMessageId(),
			Timestamp:  time.Now().UnixNano(),
			StatusCode: 500,
		}, err
	}

	tags := map[string]string{
		"deviceId":  metrics.GetDeviceId(),
		"messageId": metrics.GetMessageId(),
	}

	fields := map[string]interface{}{
		"moisture": metrics.GetMoistureLevel(),
	}

	pt, err := influx.NewPoint("moisture-level", tags, fields, time.Now())
	if err != nil {
		return &pb.MetricsReply{
			MessageId:  metrics.GetMessageId(),
			StatusCode: 500,
		}, err
	}
	bp.AddPoint(pt)
	if err := s.influxClient.WriteCtx(ctx, bp); err != nil {
		return &pb.MetricsReply{
			MessageId:  metrics.GetMessageId(),
			StatusCode: 500,
		}, err
	}

	return &pb.MetricsReply{
		MessageId:  metrics.GetMessageId(),
		StatusCode: 200,
	}, nil
}
