package client

const (
	// StatusLive is the game status when the game is currently live
	StatusLive = "In Progress"

	// StatusPreGame is the pre-game status, or the status during period breaks
	StatusPreGame = "Pre-Game"

	// StatusScheduled is the game status when the game hasn't started yet
	StatusScheduled = "Scheduled"

	// StatusFinal is the game status once the game is finished
	StatusFinal = "Final"
)

// ScheduleDate represents a date entry in the schedule response
type ScheduleDate struct {
	Date  string          `json:"date"`
	Games []*ScheduleGame `json:"games"`
}

// ScheduleGame is a scheduled game for a team
type ScheduleGame struct {
	ID       int    `json:"gamePk"`
	GameDate string `json:"gameDate"`
	Status   struct {
		Detailed string `json:"detailedState"`
	} `json:"status"`
	Teams struct {
		Home HomeAwayTeam `json:"home"`
		Away HomeAwayTeam `json:"away"`
	} `json:"teams"`
	Link string `json:"link"`
}

// HomeAwayTeam represents the Home or Away team
type HomeAwayTeam struct {
	Team  ScheduleGameTeam `json:"team"`
	Score int              `json:"score"`
}

// ScheduleGameTeam contains the team ID for a game
type ScheduleGameTeam struct {
	ID int `json:"id"`
}

// ScheduleResponse is returned from the NHL api for schedule information
type ScheduleResponse struct {
	Dates []*ScheduleDate `json:"dates"`
}
