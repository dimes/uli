package client

const (
	// CurrentPeriodEnded is the value of CurrentPeriodTimeRemaining when the period is over
	CurrentPeriodEnded = "END"
)

// LiveGame is returned from a game link when it is live
type LiveGame struct {
	LiveData struct {
		Linescore struct {
			CurrentPeriod              int    `json:"currentPeriod"`
			CurrentPeriodTimeRemaining string `json:"currentPeriodTimeRemaining"`
		} `json:"linescore"`
	} `json:"liveData"`
}
