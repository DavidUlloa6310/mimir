package models

// Ticket represents a ticket structure
type Ticket struct {
	Number           int    `json:"number"`
	ShortDescription string `json:"short_description"`
	State            string `json:"state"`
	Priority         string `json:"priority"`
}
