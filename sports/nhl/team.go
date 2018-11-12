package nhl

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	nhl_client "github.com/dimes/uli/sports/nhl/client"
	"github.com/pkg/errors"
)

// Team are commands for a single team
type Team struct {
	ID           int
	Abbreviation string

	client    nhl_client.Client
	teamCache *TeamCache
}

// NewTeam returns a new team command
func NewTeam(id int, abbreviation string, client nhl_client.Client, teamCache *TeamCache) *Team {
	return &Team{
		ID:           id,
		Abbreviation: abbreviation,
		client:       client,
		teamCache:    teamCache,
	}
}

// FromTeam creates a team command from the API response
func FromTeam(client nhl_client.Client, teamCache *TeamCache, team *nhl_client.Team) *Team {
	return NewTeam(team.ID, team.Abbreviation, client, teamCache)
}

// Name returns the name of the team command
func (t *Team) Name() string {
	return t.Abbreviation
}

// Execute executes the team command
func (t *Team) Execute(ctx context.Context, args []string) error {
	team, err := t.teamCache.GetTeam(ctx, t.ID)
	if err != nil {
		return errors.Wrapf(err, "Error getting team for id %d", t.ID)
	}

	if team == nil {
		return fmt.Errorf("No team found for id %d", t.ID)
	}

	oneWeekAgo := time.Now().Add(-7 * 24 * time.Hour).Local().Format("2006-01-02")
	oneWeekFromNow := time.Now().Add(7 * 24 * time.Hour).Local().Format("2006-01-02")
	games, err := t.client.Schedule(ctx, team.ID, oneWeekAgo, oneWeekFromNow)
	if err != nil {
		return errors.Wrapf(err, "Error getting schedule for team %s", team.FullName)
	}

	sort.Slice(games, func(i, j int) bool {
		return strings.Compare(games[i].GameDate, games[j].GameDate) < 0
	})

	return PrintGames(ctx, t.teamCache, games)
}
