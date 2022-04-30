package main

import (
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/5imonGustafsson/grpc-iot-concept/server/pb/messages"
	"github.com/5imonGustafsson/grpc-iot-concept/server/service"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"google.golang.org/grpc"
)

var (
	port        = flag.Int("port", 50051, "The server port")
	influxAddr  = flag.String("influx-address", "http://localhost:8086", "influxDB endpoint")
	influxOrg   = flag.String("influx-org", "demo.net", "influxDB organisation")
	influxToken = flag.String("influx-token", "admin", "influxDB token")
)

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	client := influxdb2.NewClient(*influxAddr, *influxToken)
	if err != nil {
		log.Fatalf("failed to setup influxDB client. Error: %v", err)
	}

	defer client.Close()

	pb.RegisterIoTServer(s, service.New(client, *influxOrg))

	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
