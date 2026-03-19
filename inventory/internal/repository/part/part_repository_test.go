package part

import (
	"testing"
	"time"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
)

func setupStorage(storage *InventoryStorage) {
	parts := []*model.Part{
		{
			UUID:          "engine-001",
			Name:          "V8 Engine",
			Description:   "High performance V8 engine",
			Price:         15000.0,
			StockQuantity: 5,
			Category:      model.CategoryEngine,
			Dimensions:    &model.Dimensions{Lenght: 100, Width: 80, Height: 60, Weight: 200.5},
			Manufacturer:  &model.Manufacturer{Name: "Ford", Country: "USA", Website: "ford.com"},
			Tags:          []string{"performance", "v8", "american"},
			Metadata:      map[string]*model.Value{"horsepower": {Int64Value: 450}},
			CreatedAt:     &time.Time{},
			UpdatedAt:     &time.Time{},
		},
		{
			UUID:          "fuel-001",
			Name:          "Fuel Pump",
			Description:   "Electric fuel pump",
			Price:         250.0,
			StockQuantity: 20,
			Category:      model.CategoryFuel,
			Dimensions:    &model.Dimensions{Lenght: 15, Width: 10, Height: 8, Weight: 0.5},
			Manufacturer:  &model.Manufacturer{Name: "Bosch", Country: "Germany", Website: "bosch.com"},
			Tags:          []string{"electric", "fuel", "german"},
			Metadata:      map[string]*model.Value{"flow_rate": {DoubleValue: 100.5}},
			CreatedAt:     &time.Time{},
			UpdatedAt:     &time.Time{},
		},
		{
			UUID:          "wing-001",
			Name:          "Aircraft Wing",
			Description:   "Commercial aircraft wing",
			Price:         50000.0,
			StockQuantity: 2,
			Category:      model.CategoryWing,
			Dimensions:    &model.Dimensions{Lenght: 2000, Width: 300, Height: 50, Weight: 1500.0},
			Manufacturer:  &model.Manufacturer{Name: "Boeing", Country: "USA", Website: "boeing.com"},
			Tags:          []string{"aircraft", "commercial", "large"},
			Metadata:      map[string]*model.Value{"wing_span": {DoubleValue: 60.0}},
			CreatedAt:     &time.Time{},
			UpdatedAt:     &time.Time{},
		},
		{
			UUID:          "porthole-001",
			Name:          "Round Porthole",
			Description:   "Standard round porthole",
			Price:         800.0,
			StockQuantity: 15,
			Category:      model.CategoryPorthole,
			Dimensions:    &model.Dimensions{Lenght: 30, Width: 30, Height: 5, Weight: 2.5},
			Manufacturer:  &model.Manufacturer{Name: "MarineCo", Country: "Italy", Website: "marineco.it"},
			Tags:          []string{"marine", "round", "standard"},
			Metadata:      map[string]*model.Value{"diameter": {Int64Value: 30}},
			CreatedAt:     &time.Time{},
			UpdatedAt:     &time.Time{},
		},
		{
			UUID:          "engine-002",
			Name:          "Turbo Engine",
			Description:   "Turbocharged 4-cylinder engine",
			Price:         8000.0,
			StockQuantity: 8,
			Category:      model.CategoryEngine,
			Dimensions:    &model.Dimensions{Lenght: 60, Width: 50, Height: 40, Weight: 120.0},
			Manufacturer:  &model.Manufacturer{Name: "Toyota", Country: "Japan", Website: "toyota.com"},
			Tags:          []string{"turbo", "japanese", "efficient"},
			Metadata:      map[string]*model.Value{"turbo_boost": {DoubleValue: 1.5}},
			CreatedAt:     &time.Time{},
			UpdatedAt:     &time.Time{},
		},
	}
	for _, part := range parts {
		storage.Save(part)
	}
}

func TestListParts(t *testing.T) {
	repo := NewRepository()
	setupStorage(repo.GetStorage())

	testCases := []struct {
		name          string
		pf            *model.PartFilter
		expectedCount int
		expectedUuids []string
		expectError   bool
		err           error
	}{
		{
			name:          "no filter - return all parts",
			pf:            &model.PartFilter{},
			expectedCount: 5,
			expectedUuids: []string{"engine-001", "fuel-001", "wing-001", "porthole-001", "engine-002"},
			expectError:   false,
		},
		{
			name:          "nil filter - return all parts",
			pf:            nil,
			expectedCount: 5,
			expectedUuids: []string{"engine-001", "fuel-001", "wing-001", "porthole-001", "engine-002"},
			expectError:   false,
		},
		{
			name:          "filter by single UUID",
			pf:            &model.PartFilter{Uuids: []string{"engine-001"}},
			expectedCount: 1,
			expectedUuids: []string{"engine-001"},
			expectError:   false,
		},
		{
			name:          "filter by multiple UUIDs",
			pf:            &model.PartFilter{Uuids: []string{"engine-001", "fuel-001"}},
			expectedCount: 2,
			expectedUuids: []string{"engine-001", "fuel-001"},
			expectError:   false,
		},
		{
			name:          "filter by non-existent UUID",
			pf:            &model.PartFilter{Uuids: []string{"non-existent"}},
			expectedCount: 0,
			expectedUuids: []string{},
			expectError:   false,
		},
		{
			name:          "filter by single name",
			pf:            &model.PartFilter{Names: []string{"V8 Engine"}},
			expectedCount: 1,
			expectedUuids: []string{"engine-001"},
			expectError:   false,
		},
		{
			name:          "filter by multiple names",
			pf:            &model.PartFilter{Names: []string{"V8 Engine", "Fuel Pump"}},
			expectedCount: 2,
			expectedUuids: []string{"engine-001", "fuel-001"},
			expectError:   false,
		},
		{
			name:          "filter by non-existent name",
			pf:            &model.PartFilter{Names: []string{"Non-existent Part"}},
			expectedCount: 0,
			expectedUuids: []string{},
			expectError:   false,
		},
		{
			name:          "filter by single category",
			pf:            &model.PartFilter{Categories: []model.Category{model.CategoryEngine}},
			expectedCount: 2,
			expectedUuids: []string{"engine-001", "engine-002"},
			expectError:   false,
		},
		{
			name:          "filter by multiple categories",
			pf:            &model.PartFilter{Categories: []model.Category{model.CategoryEngine, model.CategoryFuel}},
			expectedCount: 3,
			expectedUuids: []string{"engine-001", "fuel-001", "engine-002"},
			expectError:   false,
		},
		{
			name:          "filter by single manufacturer country",
			pf:            &model.PartFilter{ManufacturerCountries: []string{"USA"}},
			expectedCount: 2,
			expectedUuids: []string{"engine-001", "wing-001"},
			expectError:   false,
		},
		{
			name:          "filter by multiple manufacturer countries",
			pf:            &model.PartFilter{ManufacturerCountries: []string{"USA", "Germany"}},
			expectedCount: 3,
			expectedUuids: []string{"engine-001", "fuel-001", "wing-001"},
			expectError:   false,
		},
		{
			name:          "filter by non-existent manufacturer country",
			pf:            &model.PartFilter{ManufacturerCountries: []string{"Non-existent"}},
			expectedCount: 0,
			expectedUuids: []string{},
			expectError:   false,
		},
		{
			name:          "filter by single tag",
			pf:            &model.PartFilter{Tags: []string{"performance"}},
			expectedCount: 1,
			expectedUuids: []string{"engine-001"},
			expectError:   false,
		},
		{
			name:          "filter by multiple tags",
			pf:            &model.PartFilter{Tags: []string{"american", "german"}},
			expectedCount: 2,
			expectedUuids: []string{"engine-001", "fuel-001"},
			expectError:   false,
		},
		{
			name:          "filter by non-existent tag",
			pf:            &model.PartFilter{Tags: []string{"non-existent"}},
			expectedCount: 0,
			expectedUuids: []string{},
			expectError:   false,
		},
		{
			name:          "complex filter - UUID + category",
			pf:            &model.PartFilter{Uuids: []string{"engine-001"}, Categories: []model.Category{model.CategoryEngine}},
			expectedCount: 1,
			expectedUuids: []string{"engine-001"},
			expectError:   false,
		},
		{
			name: "complex filter - name + country + tags",
			pf: &model.PartFilter{
				Names:                 []string{"V8 Engine"},
				ManufacturerCountries: []string{"USA"},
				Tags:                  []string{"performance"},
			},
			expectedCount: 1,
			expectedUuids: []string{"engine-001"},
			expectError:   false,
		},
		{
			name: "complex filter - conflicting criteria",
			pf: &model.PartFilter{
				Uuids:      []string{"engine-001"},
				Categories: []model.Category{model.CategoryFuel},
			},
			expectedCount: 0,
			expectedUuids: []string{},
			expectError:   false,
		},
		{
			name: "empty filter arrays",
			pf: &model.PartFilter{
				Uuids:                 []string{},
				Names:                 []string{},
				Categories:            []model.Category{},
				ManufacturerCountries: []string{},
				Tags:                  []string{},
			},
			expectedCount: 5,
			expectedUuids: []string{"engine-001", "fuel-001", "wing-001", "porthole-001", "engine-002"},
			expectError:   false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := repo.ListParts(t.Context(), tc.pf)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tc.err != nil && err.Error() != tc.err.Error() {
					t.Errorf("expected error %v, got %v", tc.err, err)
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if len(result) != tc.expectedCount {
				t.Errorf("expected %d parts, got %d", tc.expectedCount, len(result))
			}

			if len(tc.expectedUuids) > 0 {
				resultUuids := make([]string, len(result))
				for i, part := range result {
					resultUuids[i] = part.GetUUID()
				}

				for _, expectedUUID := range tc.expectedUuids {
					found := false
					for _, resultUUID := range resultUuids {
						if resultUUID == expectedUUID {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("expected UUID %s not found in result", expectedUUID)
					}
				}
			}
		})
	}
}
