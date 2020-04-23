package models

import (
	"errors"
	"strings"
)

func (p *Product) Validate() error {

	// Validate the sort by field

	var listErr []string

	if strings.Trim(p.Name, " ") == "" {
		listErr = append(listErr, "| Name Required |")
	}

	if strings.Trim(p.Description, " ") == "" {
		listErr = append(listErr, "| Description Required |")
	}

	if p.Price <= 0 {
		listErr = append(listErr, "| Price value is required |")
	}

	if p.DeliveryPrice < 0 {
		listErr = append(listErr, "| Invalid Delivery price |")
	}

	if len(listErr) != 0 {
		return errors.New(strings.Join(listErr, ", "))
	}

	return nil
}

func (p *ProductOption) Validate() error {

	// Validate the sort by field

	var listErr []string

	if strings.Trim(p.Name, " ") == "" {
		listErr = append(listErr, "| Name Required |")
	}

	if strings.Trim(p.Description, " ") == "" {
		listErr = append(listErr, "| Description Required |")
	}

	if len(listErr) != 0 {
		return errors.New(strings.Join(listErr, ", "))
	}

	return nil
}
