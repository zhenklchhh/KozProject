package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"buf.build/go/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"

	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type InventoryService struct {
	inventoryV1.UnimplementedInventoryServiceServer
	InventoryStorage *InventoryStorage
}

type InventoryStorage struct {
	mu      sync.RWMutex
	storage map[string]*inventoryV1.Part
}

func (s *InventoryStorage) GetAll() []*inventoryV1.Part {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]*inventoryV1.Part, 0, len(s.storage))
	for _, v := range s.storage {
		values = append(values, v)
	}
	return values
}

func (s *InventoryStorage) Get(id string) (*inventoryV1.Part, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.storage[id]
	return part, ok
}

func (s *InventoryStorage) Save(part *inventoryV1.Part) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[part.GetUuid()] = part
}

func NewStorage() *InventoryStorage {
	part1UUID := "a1b2c3d4-e5f6-7890-1234-567890abcdef"
	part2UUID := "f0e9d8c7-b6a5-4321-fedc-ba9876543210"

	storage := make(map[string]*inventoryV1.Part)
	storage[part1UUID] = &inventoryV1.Part{
		Uuid:  part1UUID,
		Name:  "GeForce RTX 4090",
		Price: 1599.99,
	}
	storage[part2UUID] = &inventoryV1.Part{
		Uuid:  part2UUID,
		Name:  "Intel Core i9-13900K",
		Price: 589.00,
	}
	return &InventoryStorage{
		storage: storage,
	}
}

func (s *InventoryService) GetPart(ctx context.Context,
	req *inventoryV1.GetPartRequest,
) (*inventoryV1.GetPartResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "inventory service: validation error")
	}
	v, ok := s.InventoryStorage.Get(req.GetUuid())
	if !ok {
		return nil, status.Errorf(codes.NotFound, "inventory service: part %s not found ", req.Uuid)
	}
	return &inventoryV1.GetPartResponse{Part: v}, nil
}

func (s *InventoryService) ListParts(ctx context.Context,
	req *inventoryV1.ListPartsRequest,
) (*inventoryV1.ListPartsResponse, error) {
	if err := protovalidate.Validate(req); err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "inventory service: validation error")
	}
	pf := req.GetFilter()
	var result []*inventoryV1.Part
	for _, part := range s.InventoryStorage.GetAll() {
		if len(pf.GetUuids()) > 0 && !contains(pf.GetUuids(), part.GetUuid()) {
			continue
		}
		if len(pf.GetNames()) > 0 && !contains(pf.GetNames(), part.GetName()) {
			continue
		}
		if len(pf.GetCategories()) > 0 && !contains(pf.GetCategories(), part.GetCategory()) {
			continue
		}
		if len(pf.GetManufacturerCountries()) > 0 && !contains(pf.GetManufacturerCountries(), part.GetManufacturer().GetCountry()) {
			continue
		}
		if len(pf.GetTags()) > 0 {
			match := false
			for _, tag := range part.GetTags() {
				if contains(pf.GetTags(), tag) {
					match = true
					break
				}
			}
			if !match {
				continue
			}
		}
		result = append(result, part)
	}
	return &inventoryV1.ListPartsResponse{Parts: result}, nil
}

func contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("failed to listen: %v\n", err)
		return
	}
	defer func() {
		if cerr := lis.Close(); cerr != nil {
			log.Printf("failed to close listener: %v\n", err)
		}
	}()
	s := grpc.NewServer()
	service := &InventoryService{
		InventoryStorage: NewStorage(),
	}
	inventoryV1.RegisterInventoryServiceServer(s, service)
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
