package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"

	"github.com/ducthangng/GeoFleet/app/internal/registry"
	"github.com/ducthangng/GeoFleet/singleton"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

func main() {
	singleton.InitializeConfig()
	centralConfig := singleton.GetGlobalConfig()

	ctx := context.Background()
	dbConn := singleton.ConnectPostgre(ctx)

	// 1. Mở port TCP
	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%s", centralConfig.Host, centralConfig.Port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// 2. Khởi tạo gRPC Server
	server := grpc.NewServer()

	// 4. Xử lý Graceful Shutdown (Hủy đăng ký khi tắt app)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt, syscall.SIGTERM)
		<-c
		log.Println("Đang tắt Service...")
		// Hủy đăng ký trên Consul để Gateway không gọi vào nữa
		// (Thực hiện gọi client.Agent().ServiceDeregister(serviceID))
		server.GracefulStop()
		os.Exit(0)
	}()

	// 1. Khởi tạo Health Server
	healthServer := health.NewServer()

	// 2. Đăng ký Health Service vào gRPC Server của bạn
	healthpb.RegisterHealthServer(server, healthServer)

	// 3. Đặt trạng thái là SERVING (Đang hoạt động)
	// "user-service" ở đây phải khớp với Service Name bạn đăng ký trên Consul
	healthServer.SetServingStatus("user-service", healthpb.HealthCheckResponse_SERVING)

	reflection.Register(server)

	// periodic health check
	go func() {
		var errDB error
		// var
		for {

			if errDB = dbConn.DB.Ping(ctx); err == nil {
				healthServer.SetServingStatus("geofleet.tracking.v1", healthpb.HealthCheckResponse_SERVING)
			} else {
				healthServer.SetServingStatus("geofleet.tracking.v1", healthpb.HealthCheckResponse_NOT_SERVING)
				log.Fatal(errDB)
			}

			time.Sleep(5 * time.Second)
		}
	}()

	// start kafka listener
	consumer, err := registry.BuildKafkaListener(ctx)
	if err != nil {
		log.Fatalf("error occur when building kafka listener: %s", err.Error())
	}

	go consumer.Listen(ctx)

	log.Println("starting server...")
	// 5. Start server
	if err := server.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
