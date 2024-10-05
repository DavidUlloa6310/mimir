package models

type Accelerator struct {
    ID    int    `json:"iD"`
    Url   string `json:"url"`
    Title string `json:"title"`
	Description string `json:"description"`
    Category string `json:"category"`
}