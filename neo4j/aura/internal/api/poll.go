package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/neo4j/cli/common/clicfg"
)

type PollResponse struct {
	Data struct {
		Id     string
		Status string
	}
}

func PollInstance(cfg *clicfg.Config, instanceId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s", instanceId)
	return poll(cfg, path, waitingStatus)
}

func PollSnapshot(cfg *clicfg.Config, instanceId string, snapshotId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s/snapshots/%s", instanceId, snapshotId)
	return poll(cfg, path, waitingStatus)
}

func PollCMK(cfg *clicfg.Config, cmkId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/customer-managed-keys/%s", cmkId)
	return poll(cfg, path, waitingStatus)
}

func poll(cfg *clicfg.Config, url string, waitingStatus string) (*PollResponse, error) {
	pollingConfig := cfg.Aura.PollingConfig()

	for i := 0; i < pollingConfig.MaxRetries; i++ {
		resBody, statusCode, err := MakeRequest(cfg, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		if statusCode == http.StatusOK {
			var response PollResponse
			if err := json.Unmarshal(resBody, &response); err != nil {
				return nil, err
			}

			if response.Data.Status == "" || response.Data.Status == waitingStatus {
				time.Sleep(time.Second * time.Duration(pollingConfig.Interval))
			} else {
				return &response, nil
			}
		} else {
			// Edge case of a status code 2xx is returned different of 200
			time.Sleep(time.Second * time.Duration(pollingConfig.Interval))
		}
	}

	return nil, fmt.Errorf("hit max retries for polling")
}
