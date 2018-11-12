package nhl

import (
	"context"
	"fmt"
	"os"
	"time"

	nhl_client "github.com/dimes/uli/sports/nhl/client"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
)

// PrintGames prints out the given list of games in a table
func PrintGames(ctx context.Context, teamCache *TeamCache, games []*nhl_client.ScheduleGame) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "", "Home", "Away", "score"})
	table.SetRowSeparator("-")
	nowLocal := time.Now().Local()
	for _, game := range games {
		gameTime, err := time.Parse(time.RFC3339, game.GameDate)
		if err != nil {
			return errors.Wrapf(err, "Error parsing game date %s", game.GameDate)
		}

		homeID := game.Teams.Home.Team.ID
		home, err := teamCache.GetTeam(ctx, homeID)
		if err != nil {
			return errors.Wrapf(err, "Error getting team %d", homeID)
		}

		if home == nil {
			return fmt.Errorf("No team with id %d found", homeID)
		}

		awayID := game.Teams.Away.Team.ID
		away, err := teamCache.GetTeam(ctx, awayID)
		if err != nil {
			return errors.Wrapf(err, "Error getting team %d", awayID)
		}

		if away == nil {
			return fmt.Errorf("No team with id %d found", awayID)
		}

		score := ""
		if showScoreStatus[game.Status.Detailed] {
			score = fmt.Sprintf("%d - %d", game.Teams.Home.Score, game.Teams.Away.Score)
		}

		gameLocal := gameTime.Local()
		formattedDate := gameLocal.Format("01/02")
		if gameLocal.Year() == nowLocal.Year() {
			if gameLocal.YearDay() == nowLocal.YearDay()-1 {
				formattedDate = "Yesterday"
			} else if gameLocal.YearDay() == nowLocal.YearDay() {
				formattedDate = "Today"
			} else if gameLocal.YearDay() == nowLocal.YearDay()+1 {
				formattedDate = "Tomorrow"
			}
		}

		formattedTime := fmt.Sprintf("%s %s", formattedDate, gameLocal.Format("3:04 pm"))
		table.Append([]string{
			formattedTime,
			statusToEmoji[game.Status.Detailed],
			home.FullName,
			away.FullName,
			score,
		})
	}
	table.Render()

	return nil
}
