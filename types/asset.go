package types

type PricingUnitEnum int

const (
	Day PricingUnitEnum = iota
	Week
	Month
	Semester
)

type ConditionEnum int

const (
	New ConditionEnum = iota
	Good
	Fair
	Poor
	Damaged
)

type AvailabilityEnum int

const (
	Available AvailabilityEnum = iota
	Paused
)

type Asset struct {
	Rate         int              `json:"rate"`
	Name         string           `json:"name"`
	Vendor       int              `json:"vendor"`
	Category     int              `json:"category"`
	Location     string           `json:"location"`
	Description  string           `json:"description"`
	PricingUnit  PricingUnitEnum  `json:"pricingUnit"`
	PrimaryImage int              `json:"primaryImage"`
	Availability AvailabilityEnum `json:"availability"`
	Condition    ConditionEnum    `json:"conditionEnum"`
}

type AssetImage struct {
	FileBytes   []byte
	ContentType string
}
