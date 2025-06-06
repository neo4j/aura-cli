package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/neo4j/cli/common/clicfg"
	"github.com/neo4j/cli/common/clierr"
)

type PollResponse struct {
	Data struct {
		Id     string
		Status string
	}
}

func PollInstance(cfg *clicfg.Config, instanceId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s", instanceId)
	return Poll(cfg, path, func(status string) bool {
		return status != waitingStatus
	})
}

func PollSnapshot(cfg *clicfg.Config, instanceId string, snapshotId string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s/snapshots/%s", instanceId, snapshotId)
	return Poll(cfg, path, func(status string) bool {
		return status != SnapshotStatusPending && status != SnapshotStatusInProgress
	})
}

func PollCMK(cfg *clicfg.Config, cmkId string) (*PollResponse, error) {
	path := fmt.Sprintf("/customer-managed-keys/%s", cmkId)
	return Poll(cfg, path, func(status string) bool {
		return status != CMKStatusPending
	})
}

func PollGraphQLDataApi(cfg *clicfg.Config, instanceId string, graphQLDataApiId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s/data-apis/graphql/%s", instanceId, graphQLDataApiId)
	return Poll(cfg, path, func(status string) bool {
		return status != waitingStatus
	})
}

func PollGraphAnalyticsSessionReady(cfg *clicfg.Config, sessionId string, waitingStatus []string) (*PollResponse, error) {
	path := fmt.Sprintf("/graph-analytics/sessions/%s", sessionId)
	return Poll(cfg, path, func(status string) bool {
		return !slices.Contains(waitingStatus, status)
	})
}

func Poll(cfg *clicfg.Config, url string, cond func(status string) bool) (*PollResponse, error) {
	pollingConfig := cfg.Aura.PollingConfig()
	for i := 0; i < pollingConfig.MaxRetries; i++ {
		time.Sleep(time.Second * time.Duration(pollingConfig.Interval))
		resBody, statusCode, err := MakeRequest(cfg, url, &RequestConfig{
			Method: http.MethodGet,
		})
		if err != nil {
			return nil, clierr.NewUpstreamError("error polling: %w", err)
		}

		if statusCode == http.StatusOK {
			var response PollResponse
			if err := json.Unmarshal(resBody, &response); err != nil {
				return nil, clierr.NewUpstreamError("cannot retrieve response polling: %w", err)
			}

			// Successful poll, return last response
			if cond(response.Data.Status) {
				return &response, nil
			}
		}
	}

	return nil, clierr.NewUpstreamError("hit max retries [%d] polling", pollingConfig.MaxRetries)
}
