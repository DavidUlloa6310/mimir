package models

type Accelerator struct {
    ID    int    `json:"id"`
    Url   string `json:"url"`
    Title string `json:"title"`
	Description string `json:"description"`
}