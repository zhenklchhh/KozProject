package model

import "time"

type Category int32

const (
	CategoryUnspecified Category = iota
	CategoryEngine
	CategoryFuel
	CategoryPorthole
	CategoryWing
)

var (
	Category_name = map[Category]string{
		CategoryUnspecified: "CATEGORY_UNSPECIFIED",
		CategoryEngine:      "CATEGORY_ENGINE",
		CategoryFuel:        "CATEGORY_FUEL",
		CategoryPorthole:    "CATEGORY_PORTHOLE",
		CategoryWing:        "CATEGORY_WING",
	}
	Category_value = map[string]Category{
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
	Name    string
	Country string
	Website string
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
	Uuid          string
	Name          string
	Description   string
	Price         float64
	StockQuantity int64
	Category      Category
	Dimensions    *Dimensions
	Manufacturer  *Manufacturer
	Tags          []string
	Metadata      map[string]*Value
	CreatedAt     *time.Time
	UpdatedAt     *time.Time
}

func (p *Part) GetUuid() string {
	if p.Uuid != "" {
		return p.Uuid
	}
	return ""
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
	Uuids                 []string
	Names                 []string
	Categories            []Category
	ManufacturerCountries []string
	Tags                  []string
}

func (x *PartFilter) GetUuids() []string {
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
