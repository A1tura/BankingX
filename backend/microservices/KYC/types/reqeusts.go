package types

import "time"

type KYCRequest struct {
	FirstName   string    `json:"firstName"`
	MiddleName  string    `json:"middleName"`
	LastName    string    `json:"lastName"`
	DateOfBirth time.Time `json:"dateOfBirth"`

	PhoneNumber string `json:"phoneNumber"`

	IdNumber string `json:"idNumber"`
	IdFront  string `json:"idFront"`
	IdBack   string `json:"idBack"`
	Selfie   string `json:"selfie"`

	Country string `json:"country"`
	State   string `json:"state"`
	City    string `json:"city"`

	Address    string `json:"address"`
	PostalCode string `json:"postalCode"`
}
