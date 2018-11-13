package nhl

import (
	"context"
	"sort"
	"strings"
	"time"

	nhl_client "github.com/dimes/uli/sports/nhl/client"
	"github.com/dimes/uli/util/command"
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

// NHL is the main command for NHL related information
type NHL struct {
	client    nhl_client.Client
	teamCache *TeamCache
	commands  *command.Set
}

// NewNHL returns an initialized NHL command
func NewNHL() *NHL {
	commands := []command.Command{}
	set, err := command.NewSet(commands)

	if err != nil {
		panic(err)
	}

	client := nhl_client.NewHTTPClient("https://statsapi.web.nhl.com")
	return &NHL{
		client:    client,
		teamCache: NewTeamCache(client),
		commands:  set,
	}
}

// Name returns the command name for the NHL command
func (n *NHL) Name() string {
	return "nhl"
}

// Execute runs the NHL command with the given arguments
func (n *NHL) Execute(ctx context.Context, args []string) error {
	if len(args) == 0 {
		return n.gamesToday(ctx)
	}

	if err := n.commands.Execute(ctx, args); err != nil && err != command.ErrCommandNotFound {
		return errors.Wrapf(err, "Error executing command %s", args[0])
	} else if err == nil {
		return nil
	}

	teams, err := n.client.Teams(ctx)
	if err != nil {
		return errors.Wrap(err, "Error fetching teams")
	}

	teamCommands := make([]command.Command, 0)
	for _, team := range teams {
		teamCommands = append(teamCommands, FromTeam(n.client, n.teamCache, team))
	}

	teamCommandsSet, err := command.NewSet(teamCommands)
	if err != nil {
		return errors.Wrap(err, "Error constructing team command set")
	}

	if err := teamCommandsSet.Execute(ctx, args); err != nil {
		return errors.Wrapf(err, "Error executing command %s", args[0])
	}

	return nil
}

func (n *NHL) gamesToday(ctx context.Context) error {
	teams, err := n.client.Teams(ctx)
	if err != nil {
		return errors.Wrap(err, "Error fetching games today")
	}

	today := time.Now().Local().Format("2006-01-02")
	seenGames := make(map[int]struct{})
	gamesToday := make([]*nhl_client.ScheduleGame, 0)
	group, groupCtx := errgroup.WithContext(ctx)
	for _, t := range teams {
		team := t // capture team
		if !team.Active {
			continue
		}

		group.Go(func() error {
			games, err := n.client.Schedule(groupCtx, team.ID, today, today)
			if err != nil {
				return errors.Wrapf(err, "Error getting schedule for team %s", team.FullName)
			}

			for _, game := range games {
				if _, ok := seenGames[game.ID]; ok {
					continue
				}
				seenGames[game.ID] = struct{}{}

				gamesToday = append(gamesToday, game)
			}

			return nil
		})
	}

	if err := group.Wait(); err != nil {
		return errors.Wrap(err, "Error getting game schedules")
	}

	sort.Slice(gamesToday, func(i, j int) bool {
		return strings.Compare(gamesToday[i].GameDate, gamesToday[j].GameDate) < 0
	})

	return PrintGames(ctx, n.client, n.teamCache, gamesToday)
}
