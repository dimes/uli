package client

// Team represents an NHL team
type Team struct {
	ID           int    `json:"id"`
	FullName     string `json:"name"`
	TeamName     string `json:"teamName"`
	Abbreviation string `json:"abbreviation"`
	ShortName    string `json:"shortName"`
	Active       bool   `json:"active"`
}

// TeamsResponse is the structure of the team response from the NHL API
type TeamsResponse struct {
	Teams []*Team `json:"teams"`
}
