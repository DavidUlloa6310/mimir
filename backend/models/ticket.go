package models

// Ticket represents a ticket structure
type Ticket struct {
	ID string `json:"id"`
	Number           string `json:"number"`
	ShortDescription string `json:"short_description"`
	State            string `json:"state"`
	Priority         string `json:"priority"`
}
