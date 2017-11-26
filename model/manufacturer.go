package model

type Manufacturer struct {
	Manufacturer string `json:"name"`
}

type Manufacturers struct {
	Manufacturers []Manufacturer `json:"manufacturers"`
}