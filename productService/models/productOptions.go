package models

type ProductOptions struct {
	Items *[]ProductOption `json:"Items"`
}

type ProductOption struct {
	ID          string `json:"Id"`
	Name        string `json:"Name"`
	Description string `json:"Description"`
}
