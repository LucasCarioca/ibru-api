package models

//Sugar model for the sugar table
type Sugar struct {
	Base
	Name            string  `json:"name" binding:"required"`
	GravityPerPound float32 `json:"gravity_per_pound" binding:"required"`
}
