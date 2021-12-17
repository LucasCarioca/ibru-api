package models

//Yeast model for the yeast table
type Yeast struct {
	Base
	Name         string `json:"name" binding:"required"`
	Brand        string `json:"brand" binding:"required"`
	Tolerance    string `json:"tolerance" binding:"required"`
	Attenuation  string `json:"attenuation" binding:"required"`
	Flocculation string `json:"flocculation" binding:"required"`
	Temperature  string `json:"temperature" binding:"required"`
}
