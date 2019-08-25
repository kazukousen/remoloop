package api

const ResourceAppliance Resource = "/1/appliances"

// Appliance ...
type Appliance struct {
	Type string `json:"type"`
}
