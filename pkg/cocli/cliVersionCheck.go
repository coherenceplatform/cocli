package cocli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CoherenceMetadata struct {
	PusherAppCluster      string `json:"PUSHER_APP_CLUSTER"`
	PusherAppKey          string `json:"PUSHER_APP_KEY"`
	RudderStackJsWriteKey string `json:"RUDDERSTACK_JS_WRITE_KEY"`
	CliApiVersion         string `json:"CLI_API_VERSION,omitempty"`
}

func RunCliVersionCheck() {
	// TODO: perform pre-flight metadata checks here
	client := &http.Client{}
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf("https://%s/api/public/v1/metadata", GetCoherenceDomain()),
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	res, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	defer res.Body.Close()
	bodyBytes, err := io.ReadAll(res.Body)

	metadata := &CoherenceMetadata{}
	json.Unmarshal(bodyBytes, &metadata)
	if GetCliVersion() != metadata.CliApiVersion {
		fmt.Print("WARNING: There is a newer version of cocli available. Some commands may not work as expected until you update your cocli version\n\n")
	}

	return
}
