package nhl

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	nhl_client "github.com/dimes/uli/sports/nhl/client"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
)

var (
	statusToEmoji = map[string]string{
		nhl_client.StatusLive:  "🔴",
		nhl_client.StatusFinal: "🏁",
	}

	showScoreStatus = map[string]bool{
		nhl_client.StatusLive:  true,
		nhl_client.StatusFinal: true,
	}
)

// PrintGames prints out the given list of games in a table
func PrintGames(
	ctx context.Context,
	client nhl_client.Client,
	teamCache *TeamCache,
	games []*nhl_client.ScheduleGame,
) error {
	var lock sync.Mutex
	gameToTime := make(map[int]string)
	group, groupCtx := errgroup.WithContext(ctx)
	for _, g := range games {
		game := g
		if game.Status.Detailed == nhl_client.StatusLive {
			group.Go(func() error {
				details, err := client.LiveGameDetails(groupCtx, game)
				if err != nil {
					return err
				}

				currentPeriod := details.LiveData.Linescore.CurrentPeriod
				timeRemaining := details.LiveData.Linescore.CurrentPeriodTimeRemaining

				lock.Lock()
				defer lock.Unlock()
				if timeRemaining == nhl_client.CurrentPeriodEnded {
					gameToTime[game.ID] = fmt.Sprintf("End of P%d", currentPeriod)
				} else {
					gameToTime[game.ID] = fmt.Sprintf("P%d %s", currentPeriod, timeRemaining)
				}

				return nil
			})
		}
	}

	if err := group.Wait(); err != nil {
		return errors.Wrap(err, "Error getting live game information")
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Time", "", "Home", "Away", "score", "Details"})
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
		details := gameToTime[game.ID]
		if details == "" && game.Status.Detailed == nhl_client.StatusFinal {
			details = "FINAL"
		}

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
			details,
		})
	}
	table.Render()

	return nil
}
