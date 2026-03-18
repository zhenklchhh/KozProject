package part

import (
	"context"
	"fmt"
	"sync"

	invModel "github.com/zhenklchhh/KozProject/inventory/internal/model"
)

type repository struct {
	invStorage *InventoryStorage
}

func NewRepository() *repository {
	return &repository{
		invStorage: NewStorage(),
	}
}

func (r *repository) GetStorage() *InventoryStorage{
	return r.invStorage
}

type InventoryStorage struct {
	mu      sync.RWMutex
	storage map[string]*invModel.Part
}

func NewStorage() *InventoryStorage {
	storage := make(map[string]*invModel.Part)
	return &InventoryStorage{
		storage: storage,
	}
}

func (s *InventoryStorage) GetAll() []*invModel.Part {
	s.mu.RLock()
	defer s.mu.RUnlock()
	values := make([]*invModel.Part, 0, len(s.storage))
	for _, v := range s.storage {
		values = append(values, v)
	}
	return values
}

func (s *InventoryStorage) Get(id string) (*invModel.Part, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	part, ok := s.storage[id]
	return part, ok
}

func (s *InventoryStorage) Save(part *invModel.Part) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.storage[part.GetUuid()] = part
}

func (s *repository) GetPart(_ context.Context,
	uuid string) (*invModel.Part, error) {
	v, ok := s.invStorage.Get(uuid)
	if !ok {
		return nil, fmt.Errorf("inventory service: part %s not found ", uuid)
	}
	return v, nil
}

func (s *repository) ListParts(_ context.Context,
	pf *invModel.PartFilter) ([]*invModel.Part, error) {
	var result []*invModel.Part
	for _, part := range s.invStorage.GetAll() {
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
	return result, nil
}

func contains[T comparable](slice []T, val T) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}
