package client

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/pkg/errors"
)

// Client represents all the interactions with the NHL api
type Client interface {
	Teams(ctx context.Context) ([]*Team, error)

	// Date format is "YYYY-MM-DD"
	Schedule(ctx context.Context, teamID int, startDate, endDate string) ([]*ScheduleGame, error)
}

type httpClient struct {
	httpClient *http.Client
	endpoint   string
}

// NewHTTPClient returns an http-backed NHL client pointing at the given endpoint
func NewHTTPClient(endpoint string) Client {
	return &httpClient{
		httpClient: http.DefaultClient,
		endpoint:   endpoint,
	}
}

func (h *httpClient) Teams(ctx context.Context) ([]*Team, error) {
	res := &TeamsResponse{}
	if err := h.makeRequest(ctx, "/api/v1/teams", res); err != nil {
		return nil, errors.Wrap(err, "Error fetching teams")
	}

	return res.Teams, nil
}

func (h *httpClient) Schedule(
	ctx context.Context,
	teamID int,
	startDate string,
	endDate string,
) ([]*ScheduleGame, error) {
	res := &ScheduleResponse{}
	path := fmt.Sprintf("/api/v1/schedule?teamId=%d&startDate=%s&endDate=%s", teamID, startDate, endDate)
	if err := h.makeRequest(ctx, path, res); err != nil {
		return nil, errors.Wrap(err, "Error fetching schedule")
	}

	games := make([]*ScheduleGame, 0)
	for _, date := range res.Dates {
		games = append(games, date.Games...)
	}

	return games, nil
}

func (h *httpClient) makeRequest(ctx context.Context, path string, res interface{}) error {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s%s", h.endpoint, path), nil)
	if err != nil {
		return errors.Wrapf(err, "Error creating request to %s", path)
	}

	httpRes, err := h.httpClient.Do(req.WithContext(ctx))
	if err != nil {
		return errors.Wrapf(err, "Error making request to %s", path)
	}

	if err := json.NewDecoder(httpRes.Body).Decode(res); err != nil {
		return errors.Wrapf(err, "Error decoding response to %s", path)
	}

	return nil
}
