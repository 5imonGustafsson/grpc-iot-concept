package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	"github.com/5imonGustafsson/grpc-iot-concept/server/service"
	influx "github.com/influxdata/influxdb/client/v2"
	"google.golang.org/grpc"
)

var (
	port       = flag.Int("port", 50051, "The server port")
	influxAddr = flag.String("influx-address", "http://localhost:8086", "influxDB endpoint")
	influxUser = flag.String("influx-user", "admin", "influxDB user")
	influxPwd  = flag.String("influx-password", "admin", "influxDB password")
)

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

	pb.RegisterIoTServer(s, service.New(influxDB))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
