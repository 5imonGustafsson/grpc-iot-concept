package service

import (
	"context"
	"log"
	"time"

	pb "github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
)

// server is used to implement messages.IoT
type service struct {
	influxClient influxdb2.Client
	influxOrg    string
	pb.UnimplementedIoTServer
}

func New(influxClient influxdb2.Client, influxOrg string) pb.IoTServer {
	return &service{
		influxClient: influxClient,
		influxOrg:    influxOrg,
	}
}

func (s *service) SendWaterSoilLevel(ctx context.Context, metrics *pb.WaterSoilMetrics) (*pb.MetricsReply, error) {
	messageID := metrics.GetMessageId()
	deviceID := metrics.GetDeviceId()
	log.Printf("Received message: %v from device: %v", messageID, deviceID)

	writeAPI := s.influxClient.WriteAPIBlocking(s.influxOrg, "hydrophonic")

	// Create point using fluent style
	p := influxdb2.NewPointWithMeasurement("moisture-level").
		AddTag("device_id", deviceID).
		AddTag("message_id", messageID).
		AddField("moisture", metrics.GetMoistureLevel()).
		SetTime(time.Now())

	if err := writeAPI.WritePoint(ctx, p); err != nil {
		return &pb.MetricsReply{
			MessageId:  messageID,
			Timestamp:  time.Now().UnixNano(),
			StatusCode: 500,
		}, err
	}

	return &pb.MetricsReply{
		MessageId:  metrics.GetMessageId(),
		StatusCode: 200,
	}, nil
}
