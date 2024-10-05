package models

type Suggestion struct {
	ID           int               `json:"id"`
	Title        string            `json:"title"`
	Description  string            `json:"description"`
	Accelerators []Accelerator `json:"accelerators"` // Assuming the Accelerator is part of your models package
}