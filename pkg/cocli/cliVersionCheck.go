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
	fmt.Println(metadata)

	return
}
