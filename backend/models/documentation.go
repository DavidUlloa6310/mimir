package models

// Documentation represents a documentation entry
type Documentation struct {
	Title   string `json:"title"`
	Accelerator Accelerator `json:"accelerator"`
}