package cocli

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type CoherenceMetadata struct {
	PusherAppCluster      string `json:"PUSHER_APP_CLUSTER"`
	PusherAppKey          string `json:"PUSHER_APP_KEY"`
	RudderStackJsWriteKey string `json:"RUDDERSTACK_JS_WRITE_KEY"`
	CliApiVersion         string `json:"CLI_API_VERSION,omitempty"`
}

func RunCliVersionCheck() {
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

	parsedCliVersion, err := semverToFloat(GetCliVersion())
	if err != nil {
		panic(err)
	}
	parsedMetaCliVersion, err := semverToFloat(metadata.CliApiVersion)
	if err != nil {
		panic(err)
	}

	if parsedCliVersion < parsedMetaCliVersion {
		fmt.Print("WARNING: There is a newer version of cocli available. Some commands may not work as expected until you update your cocli version\n\n")
	}

	return
}

func semverToFloat(version string) (float64, error) {
	// Split the version string into major and minor parts using the dot as the separator
	parts := strings.Split(version, ".")

	if len(parts) < 2 {
		return 0.0, fmt.Errorf("Invalid semver version string: %s", version)
	}

	// Parse the major component into an integer
	major, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0.0, err
	}

	// Parse the minor component into a float
	minorParts := strings.Split(parts[1], "")
	if len(minorParts) == 0 {
		return 0.0, fmt.Errorf("Invalid semver version string: %s", version)
	}

	minorFloat, err := strconv.ParseFloat("0."+strings.Join(minorParts, ""), 64)
	if err != nil {
		return 0.0, err
	}

	// Combine major and minor components into a float
	versionFloat := float64(major) + minorFloat

	return versionFloat, nil
}
