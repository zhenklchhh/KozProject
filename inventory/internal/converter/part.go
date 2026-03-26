package converter

import (
	"time"

	"github.com/google/uuid"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zhenklchhh/KozProject/inventory/internal/model"
	inventoryV1 "github.com/zhenklchhh/KozProject/shared/pkg/proto/inventory/v1"
)

// PartServiceToPartRepo конвертирует protobuf Part в repository model Part.
func PartServiceToPartRepo(part *inventoryV1.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		UUID:          uuid.MustParse(part.GetUuid()),
		Name:          part.GetName(),
		Description:   part.GetDescription(),
		Price:         part.GetPrice(),
		StockQuantity: part.GetStockQuantity(),
		Category:      CategoryServiceToRepo(part.GetCategory()),
		Dimensions:    DimensionsServiceToRepo(part.GetDimensions()),
		Manufacturer:  ManufacturerServiceToRepo(part.GetManufacturer()),
		Tags:          part.GetTags(),
		Metadata:      MetadataServiceToRepo(part.GetMetadata()),
		CreatedAt:     TimestampServiceToRepo(part.GetCreatedAt()),
		UpdatedAt:     TimestampServiceToRepo(part.GetUpdatedAt()),
	}
}

// CategoryServiceToRepo конвертирует protobuf Category в repository model Category.
func CategoryServiceToRepo(category inventoryV1.Category) model.Category {
	switch category {
	case inventoryV1.Category_CATEGORY_UNSPECIFIED:
		return model.CategoryUnspecified
	case inventoryV1.Category_CATEGORY_ENGINE:
		return model.CategoryEngine
	case inventoryV1.Category_CATEGORY_FUEL:
		return model.CategoryFuel
	case inventoryV1.Category_CATEGORY_PORTHOLE:
		return model.CategoryPorthole
	case inventoryV1.Category_CATEGORY_WING:
		return model.CategoryWing
	default:
		return model.CategoryUnspecified
	}
}

// DimensionsServiceToRepo конвертирует protobuf Dimensions в repository model Dimensions.
func DimensionsServiceToRepo(dimensions *inventoryV1.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Lenght: dimensions.GetLength(),
		Width:  dimensions.GetWidth(),
		Height: dimensions.GetHeight(),
		Weight: dimensions.GetWeight(),
	}
}

// ManufacturerServiceToRepo конвертирует protobuf Manufacturer в repository model Manufacturer.
func ManufacturerServiceToRepo(manufacturer *inventoryV1.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.GetName(),
		Country: manufacturer.GetCountry(),
		Website: manufacturer.GetWebsite(),
	}
}

// ValueServiceToRepo конвертирует protobuf Value в repository model Value.
func ValueServiceToRepo(value *inventoryV1.Value) *model.Value {
	if value == nil {
		return nil
	}

	switch v := value.GetValue().(type) {
	case *inventoryV1.Value_StringValue:
		return &model.Value{StringValue: v.StringValue}
	case *inventoryV1.Value_Int64Value:
		return &model.Value{Int64Value: v.Int64Value}
	case *inventoryV1.Value_DoubleValue:
		return &model.Value{DoubleValue: v.DoubleValue}
	case *inventoryV1.Value_BoolValue:
		return &model.Value{BoolValue: v.BoolValue}
	default:
		return nil
	}
}

// MetadataServiceToRepo конвертирует protobuf metadata map в repository model metadata map.
func MetadataServiceToRepo(metadata map[string]*inventoryV1.Value) map[string]*model.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*model.Value)
	for key, value := range metadata {
		result[key] = ValueServiceToRepo(value)
	}
	return result
}

// TimestampServiceToRepo конвертирует protobuf Timestamp в time.Time.
func TimestampServiceToRepo(timestamp *timestamppb.Timestamp) *time.Time {
	if timestamp == nil {
		return nil
	}

	t := timestamp.AsTime()
	return &t
}

// PartFilterServiceToRepo конвертирует protobuf PartFilter в repository model PartFilter.
func PartFilterServiceToRepo(filter *inventoryV1.PartFilter) *model.PartFilter {
	if filter == nil {
		return nil
	}

	categories := make([]model.Category, 0, len(filter.GetCategories()))
	for _, category := range filter.GetCategories() {
		categories = append(categories, CategoryServiceToRepo(category))
	}

	uuids := make([]uuid.UUID, 0, len(filter.Uuids))
	for _, id := range filter.GetUuids() {
		uuids = append(uuids, uuid.MustParse(id))
	}

	return &model.PartFilter{
		Uuids:                 uuids,
		Names:                 filter.GetNames(),
		Categories:            categories,
		ManufacturerCountries: filter.GetManufacturerCountries(),
		Tags:                  filter.GetTags(),
	}
}

// PartRepoToService конвертирует repository model Part в protobuf Part.
func PartRepoToService(part *model.Part) *inventoryV1.Part {
	if part == nil {
		return nil
	}

	return &inventoryV1.Part{
		Uuid:          part.UUID.String(),
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryRepoToService(part.Category),
		Dimensions:    DimensionsRepoToService(part.Dimensions),
		Manufacturer:  ManufacturerRepoToService(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataRepoToService(part.Metadata),
		CreatedAt:     TimestampRepoToService(part.CreatedAt),
		UpdatedAt:     TimestampRepoToService(part.UpdatedAt),
	}
}

// CategoryRepoToService конвертирует repository model Category в protobuf Category.
func CategoryRepoToService(category model.Category) inventoryV1.Category {
	switch category {
	case model.CategoryUnspecified:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	case model.CategoryEngine:
		return inventoryV1.Category_CATEGORY_ENGINE
	case model.CategoryFuel:
		return inventoryV1.Category_CATEGORY_FUEL
	case model.CategoryPorthole:
		return inventoryV1.Category_CATEGORY_PORTHOLE
	case model.CategoryWing:
		return inventoryV1.Category_CATEGORY_WING
	default:
		return inventoryV1.Category_CATEGORY_UNSPECIFIED
	}
}

// DimensionsRepoToService конвертирует repository model Dimensions в protobuf Dimensions.
func DimensionsRepoToService(dimensions *model.Dimensions) *inventoryV1.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &inventoryV1.Dimensions{
		Length: dimensions.Lenght,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// ManufacturerRepoToService конвертирует repository model Manufacturer в protobuf Manufacturer.
func ManufacturerRepoToService(manufacturer *model.Manufacturer) *inventoryV1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryV1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ValueRepoToService конвертирует repository model Value в protobuf Value.
func ValueRepoToService(value *model.Value) *inventoryV1.Value {
	if value == nil {
		return nil
	}

	if value.StringValue != "" {
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_StringValue{StringValue: value.StringValue},
		}
	}
	if value.Int64Value != 0 {
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_Int64Value{Int64Value: value.Int64Value},
		}
	}
	if value.DoubleValue != 0 {
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_DoubleValue{DoubleValue: value.DoubleValue},
		}
	}
	if value.BoolValue {
		return &inventoryV1.Value{
			Value: &inventoryV1.Value_BoolValue{BoolValue: value.BoolValue},
		}
	}

	return nil
}

// MetadataRepoToService конвертирует repository model metadata map в protobuf metadata map.
func MetadataRepoToService(metadata map[string]*model.Value) map[string]*inventoryV1.Value {
	if metadata == nil {
		return nil
	}

	result := make(map[string]*inventoryV1.Value)
	for key, value := range metadata {
		result[key] = ValueRepoToService(value)
	}
	return result
}

// TimestampRepoToService конвертирует time.Time в protobuf Timestamp.
func TimestampRepoToService(timestamp *time.Time) *timestamppb.Timestamp {
	if timestamp == nil {
		return nil
	}

	return timestamppb.New(*timestamp)
}

// PartFilterRepoToService конвертирует repository model PartFilter в protobuf PartFilter.
func PartFilterRepoToService(filter *model.PartFilter) *inventoryV1.PartFilter {
	if filter == nil {
		return nil
	}

	categories := make([]inventoryV1.Category, 0, len(filter.Categories))
	for _, category := range filter.Categories {
		categories = append(categories, CategoryRepoToService(category))
	}

	uuids := make([]string, 0, len(filter.Uuids))
	for _, id := range filter.GetUuids() {
		uuids = append(uuids, id.String())
	}

	return &inventoryV1.PartFilter{
		Uuids:                 uuids,
		Names:                 filter.Names,
		Categories:            categories,
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}
