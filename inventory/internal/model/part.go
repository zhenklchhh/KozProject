package model

import (
	"time"

	"github.com/google/uuid"
)

type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

var (
	CategoryName = map[Category]string{
		CategoryUnspecified: "CATEGORY_UNSPECIFIED",
		CategoryEngine:      "CATEGORY_ENGINE",
		CategoryFuel:        "CATEGORY_FUEL",
		CategoryPorthole:    "CATEGORY_PORTHOLE",
		CategoryWing:        "CATEGORY_WING",
	}
	CategoryValue = map[string]Category{
		"CATEGORY_UNSPECIFIED": CategoryUnspecified,
		"CATEGORY_ENGINE":      CategoryEngine,
		"CATEGORY_FUEL":        CategoryFuel,
		"CATEGORY_PORTHOLE":    CategoryPorthole,
		"CATEGORY_WING":        CategoryWing,
	}
)

type Dimensions struct {
	Lenght float64
	Width  float64
	Height float64
	Weight float64
}

type Manufacturer struct {
	Name    string `bson:"name"`
	Country string `bson:"country"`
	Website string `bson:"website"`
}

func (m *Manufacturer) GetCountry() string {
	return m.Country
}

type Value struct {
	Int64Value  int64
	DoubleValue float64
	BoolValue   bool
	StringValue string
}

type Part struct {
	UUID          uuid.UUID         `bson:"_id"`
	Name          string            `bson:"name"`
	Description   string            `bson:"description"`
	Price         float64           `bson:"price"`
	StockQuantity int64             `bson:"stock_quantity"`
	Category      Category          `bson:"category"`
	Dimensions    *Dimensions       `bson:"dimensions"`
	Manufacturer  *Manufacturer     `bson:"manufacturer"`
	Tags          []string          `bson:"tags"`
	Metadata      map[string]*Value `bson:"metadata"`
	CreatedAt     *time.Time        `bson:"created_at"`
	UpdatedAt     *time.Time        `bson:"updated_at"`
}

func (p *Part) GetUUID() uuid.UUID {
	if p.UUID != uuid.Nil {
		return p.UUID
	}
	return uuid.Nil
}

func (p *Part) GetName() string {
	if p.Name != "" {
		return p.Name
	}
	return ""
}

func (p *Part) GetCategory() Category {
	return p.Category
}

func (p *Part) GetManufacturer() *Manufacturer {
	if p.Manufacturer != nil {
		return p.Manufacturer
	}
	return p.Manufacturer
}

func (p *Part) GetTags() []string {
	return p.Tags
}

type PartFilter struct {
	Uuids                 []uuid.UUID
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

func (x *PartFilter) GetUuids() []uuid.UUID {
	if x != nil {
		return x.Uuids
	}
	return nil
}

func (x *PartFilter) GetNames() []string {
	if x != nil {
		return x.Names
	}
	return nil
}

func (x *PartFilter) GetCategories() []Category {
	if x != nil {
		return x.Categories
	}
	return nil
}

func (x *PartFilter) GetManufacturerCountries() []string {
	if x != nil {
		return x.ManufacturerCountries
	}
	return nil
}

func (x *PartFilter) GetTags() []string {
	if x != nil {
		return x.Tags
	}
	return nil
}
