package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/cobra"
)

const maxPollRetries = 100
const pollWaitSeconds = 20

type PollResponse struct {
	Data struct {
		Id     string
		Status string
	}
}

func PollInstance(cmd *cobra.Command, instanceId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/instances/%s", instanceId)
	return poll(cmd, path, waitingStatus)
}

func PollCMK(cmd *cobra.Command, cmkId string, waitingStatus string) (*PollResponse, error) {
	path := fmt.Sprintf("/customer-managed-keys/%s", cmkId)
	return poll(cmd, path, waitingStatus)
}

func poll(cmd *cobra.Command, url string, waitingStatus string) (*PollResponse, error) {
	for i := 0; i < maxPollRetries; i++ {
		resBody, statusCode, err := MakeRequest(cmd, http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}

		if statusCode == http.StatusOK {
			var response PollResponse
			if err := json.Unmarshal(resBody, &response); err != nil {
				return nil, err
			}

			if response.Data.Status == "" || response.Data.Status == waitingStatus {
				time.Sleep(time.Second * pollWaitSeconds)
			} else {
				return &response, nil
			}
		} else {
			// Edge case of a status code 2xx is returned different of 200
			time.Sleep(time.Second * pollWaitSeconds)
		}
	}

	return nil, fmt.Errorf("hit max retries for polling")
}
