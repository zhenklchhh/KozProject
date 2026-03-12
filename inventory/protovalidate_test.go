package inventory

import (
	"errors"
	"testing"

	"buf.build/go/protovalidate"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

func TestPartValidationWithProtovalidate(t *testing.T) {
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		part    *inventoryv1.Part
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid part",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Test Part",
				Description:   "Test Description",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				Manufacturer: &inventoryv1.Manufacturer{
					Name:    "Test Manufacturer",
					Country: "Test Country",
					Website: "https://example.com",
				},
				Tags:      []string{"tag1", "tag2"},
				Metadata:  map[string]*inventoryv1.Value{},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: false,
		},
		{
			name: "invalid uuid format",
			part: &inventoryv1.Part{
				Uuid:          "invalid-uuid",
				Name:          "Test Part",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "uuid",
		},
		{
			name: "name too short",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "A",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "min_len",
		},
		{
			name: "name too long",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "This is a very long name that exceeds the maximum allowed length of fifty five characters",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "max_len",
		},
		{
			name: "negative price",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Test Part",
				Price:         -10.0,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "gt",
		},
		{
			name: "zero stock quantity",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Test Part",
				Price:         100.50,
				StockQuantity: 0,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: 10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "gt",
		},
		{
			name: "negative dimensions",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Test Part",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Dimensions: &inventoryv1.Dimensions{
					Length: -10.0,
					Width:  5.0,
					Height: 3.0,
					Weight: 1.5,
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "gt",
		},
		{
			name: "manufacturer name too short",
			part: &inventoryv1.Part{
				Uuid:          "550e8400-e29b-41d4-a716-446655440000",
				Name:          "Test Part",
				Price:         100.50,
				StockQuantity: 10,
				Category:      inventoryv1.Category_CATEGORY_ENGINE,
				Manufacturer: &inventoryv1.Manufacturer{
					Name:    "A",
					Country: "Test Country",
				},
				CreatedAt: timestamppb.Now(),
				UpdatedAt: timestamppb.Now(),
			},
			wantErr: true,
			errMsg:  "min_len",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.part)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected validation error, but got none")
					return
				}

				var validationErr *protovalidate.ValidationError
				if errors.As(err, &validationErr) {
					t.Logf("Validation error: %v", validationErr)
					if tt.errMsg != "" && !contains(validationErr.Error(), tt.errMsg) {
						t.Errorf("Expected error message to contain '%s', got: %v", tt.errMsg, validationErr.Error())
					}
				} else {
					t.Errorf("Expected ValidationError, got: %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, but got: %v", err)
				}
			}
		})
	}
}

func TestPartFilterValidationWithProtovalidate(t *testing.T) {
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		filter  *inventoryv1.PartFilter
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid filter with uuids",
			filter: &inventoryv1.PartFilter{
				Uuids: []string{"550e8400-e29b-41d4-a716-446655440000", "550e8400-e29b-41d4-a716-446655440001"},
			},
			wantErr: false,
		},
		{
			name: "valid filter with names",
			filter: &inventoryv1.PartFilter{
				Names: []string{"Engine Part", "Filter"},
			},
			wantErr: false,
		},
		{
			name: "invalid uuid in filter",
			filter: &inventoryv1.PartFilter{
				Uuids: []string{"550e8400-e29b-41d4-a716-446655440000", "invalid-uuid"},
			},
			wantErr: true,
			errMsg:  "uuid",
		},
		{
			name: "name too short in filter",
			filter: &inventoryv1.PartFilter{
				Names: []string{"Engine Part", "A"},
			},
			wantErr: true,
			errMsg:  "min_len",
		},
		{
			name: "manufacturer country too short",
			filter: &inventoryv1.PartFilter{
				ManufacturerCountries: []string{"USA", "A"},
			},
			wantErr: true,
			errMsg:  "min_len",
		},
		{
			name: "tag too short",
			filter: &inventoryv1.PartFilter{
				Tags: []string{"engine", "A"},
			},
			wantErr: true,
			errMsg:  "min_len",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.filter)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected validation error, but got none")
					return
				}

				var validationErr *protovalidate.ValidationError
				if errors.As(err, &validationErr) {
					t.Logf("Validation error: %v", validationErr)
					if tt.errMsg != "" && !contains(validationErr.Error(), tt.errMsg) {
						t.Errorf("Expected error message to contain '%s', got: %v", tt.errMsg, validationErr.Error())
					}
				} else {
					t.Errorf("Expected ValidationError, got: %T", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, but got: %v", err)
				}
			}
		})
	}
}

func TestGetPartRequestValidation(t *testing.T) {
	validator, err := protovalidate.New()
	if err != nil {
		t.Fatalf("Failed to create validator: %v", err)
	}

	tests := []struct {
		name    string
		request *inventoryv1.GetPartRequest
		wantErr bool
	}{
		{
			name: "valid request",
			request: &inventoryv1.GetPartRequest{
				Uuid: "550e8400-e29b-41d4-a716-446655440000",
			},
			wantErr: false,
		},
		{
			name: "invalid uuid",
			request: &inventoryv1.GetPartRequest{
				Uuid: "invalid-uuid",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validator.Validate(tt.request)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected validation error, but got none")
				} else {
					t.Logf("Validation error: %v", err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no validation error, but got: %v", err)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && findSubstring(s, substr)))
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
