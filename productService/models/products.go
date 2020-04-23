package models

type Products struct {
	Items *[]Product `json:"Items"`
}

type Product struct {
	ID            string  `json:"Id"`
	Name          string  `json:"Name"`
	Description   string  `json:"Description"`
	Price         float64 `json:"Price"`
	DeliveryPrice float64 `json:"DeliveryPrice"`
}
