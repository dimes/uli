package nhl

import (
	"context"
	"sync"

	nhl_client "github.com/dimes/uli/sports/nhl/client"
)

// TeamCache stores team information in a cache
type TeamCache struct {
	client              nhl_client.Client
	teams               map[int]*nhl_client.Team
	initializer         sync.Once
	initializationError error
}

// NewTeamCache returns a new team cache
func NewTeamCache(client nhl_client.Client) *TeamCache {
	return &TeamCache{
		client: client,
		teams:  make(map[int]*nhl_client.Team),
	}
}

// Teams returns a map of team id to team
func (t *TeamCache) Teams(ctx context.Context) (map[int]*nhl_client.Team, error) {
	t.initialize(ctx)
	return t.teams, t.initializationError
}

// GetTeam returns the team for the given id, or an error if the team couldn't be fetched
func (t *TeamCache) GetTeam(ctx context.Context, id int) (*nhl_client.Team, error) {
	t.initialize(ctx)
	return t.teams[id], t.initializationError
}

func (t *TeamCache) initialize(ctx context.Context) {
	t.initializer.Do(func() {
		teams, err := t.client.Teams(ctx)
		if err != nil {
			t.initializationError = err
			return
		}

		for _, team := range teams {
			t.teams[team.ID] = team
		}
	})
}
