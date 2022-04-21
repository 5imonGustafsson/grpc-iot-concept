package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	pb "github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	influx "github.com/influxdata/influxdb/client/v2"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	influxAddr = flag.String("influx-address", "http://localhost:8086", "influxDB endpoint")
	influxUser = flag.String("influx-user", "admin", "influxDB user")
	influxPwd  = flag.String("influx-password", "admin", "influxDB password")
)

// server is used to implement messages.IoT
type server struct {
	pb.IoTServer
	influxClient influx.Client
}

func (s *server) SendWaterSoilLevel(ctx context.Context, metrics *pb.WaterSoilMetrics) (*pb.MetricsReply, error) {
	log.Printf("Received message: %v from device: %v", metrics.GetMessageId(), metrics.GetDeviceId())

	bp, err := influx.NewBatchPoints(influx.BatchPointsConfig{
		Database: "hydrophonic",
	})

	if err != nil {
		return &pb.MetricsReply{
			MessageId:  metrics.GetMessageId(),
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

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	influxDB, err := influx.NewHTTPClient(influx.HTTPConfig{
		Addr:     *influxAddr,
		Username: *influxUser,
		Password: *influxPwd,
	})
	if err != nil {
		log.Fatalf("failed to setup influxDB client. Error: %v", err)
	}

	defer influxDB.Close()

	pb.RegisterIoTServer(s, &server{influxClient: influxDB})

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
