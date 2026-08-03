package types

import (
	"encoding/json"
	"fmt"
)

type PricingUnitEnum int

const (
	Day PricingUnitEnum = iota
	Week
	Month
	Semester
)

var pricingUnitMap = map[string]PricingUnitEnum{
	"day":      Day,
	"week":     Week,
	"month":    Month,
	"semester": Semester,
}

func (p *PricingUnitEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	value, ok := pricingUnitMap[s]
	if !ok {
		return fmt.Errorf("invalid pricing unit: %q", s)
	}

	*p = value
	return nil
}

func (p PricingUnitEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.String())
}

func (p PricingUnitEnum) String() string {
	return [...]string{"day", "week", "month", "semester"}[p]
}

type ConditionEnum int

const (
	BrandNew ConditionEnum = iota
	BarelyUsed
	Standard
	FairlyUsed
	Damaged
)

var conditionMap = map[string]ConditionEnum{
	"brand_new":   BrandNew,
	"barely_used": BarelyUsed,
	"standard":    Standard,
	"fairly_used": FairlyUsed,
	"damaged":     Damaged,
}

func (c *ConditionEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	value, ok := conditionMap[s]
	if !ok {
		return fmt.Errorf("invalid condition: %q", s)
	}

	*c = value
	return nil
}

func (c ConditionEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(c.String())
}

func (c ConditionEnum) String() string {
	return [...]string{
		"brand_new",
		"barely_used",
		"standard",
		"fairly_used",
		"damaged",
	}[c]
}

type AvailabilityEnum int

const (
	Available AvailabilityEnum = iota
	Paused
)

var availabilityMap = map[string]AvailabilityEnum{
	"available": Available,
	"paused":    Paused,
}

func (a *AvailabilityEnum) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}

	value, ok := availabilityMap[s]
	if !ok {
		return fmt.Errorf("invalid availability: %q", s)
	}

	*a = value
	return nil
}

func (a AvailabilityEnum) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.String())
}

func (a AvailabilityEnum) String() string {
	return [...]string{
		"available",
		"paused",
	}[a]
}

type Asset struct {
	Rate         int              `json:"rate"`
	Name         string           `json:"name"`
	Vendor       int              `json:"vendor"`
	Category     int              `json:"category"`
	Location     string           `json:"location"`
	Condition    ConditionEnum    `json:"condition"`
	Description  string           `json:"description"`
	PricingUnit  PricingUnitEnum  `json:"pricingUnit"`
	PrimaryImage int              `json:"primaryImage"`
	Availability AvailabilityEnum `json:"availability"`
}

type PureAsset struct {
	Rate         int    `json:"rate"`
	Name         string `json:"name"`
	Vendor       int    `json:"vendor"`
	Category     string `json:"category"`
	Location     string `json:"location"`
	Condition    string `json:"condition"`
	PricingUnit  string `json:"pricingUnit"`
	PrimaryImage string `json:"primaryImage"`
}

type AssetImage struct {
	FileBytes   []byte
	ContentType string
}
