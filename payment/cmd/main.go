package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	paymentV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/payment/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const grpcPort = 50052

type PaymentService struct {
	paymentV1.UnimplementedPaymentServiceServer
}

func (s *PaymentService) PayOrder(ctx context.Context, req *paymentV1.PayOrderRequest) (*paymentV1.PayOrderResponse, error) {
	transactionUuid := uuid.New()
	log.Printf("Оплата прошла успешно, transaction_uuid: %v\n", transactionUuid)
	return &paymentV1.PayOrderResponse{TransactionUuid: transactionUuid.String()}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed listening: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", cerr)
		}
	}()
	s := grpc.NewServer()
	service := &PaymentService{}
	paymentV1.RegisterPaymentServiceServer(s, service)
	reflection.Register(s)
	go func() {
		err := s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve server: %v\n", err)
			return
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	s.GracefulStop()
}
