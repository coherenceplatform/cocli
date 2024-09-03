package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

type CocliConfig struct {
	CoherenceDomain string
}

var devConfig = &CocliConfig{
	CoherenceDomain: "https://main.c269bac2ae15.coherence.cncsites.com",
}

var prodConfig = &CocliConfig{
	CoherenceDomain: "https://beta.withcoherence.com",
}

const cliVersion = "1.1.1"

func GetCliConfig() CocliConfig {
	if strings.ToLower(os.Getenv("COHERENCE_ENVIRONMENT")) == "review" {
		return *devConfig
	}

	return *prodConfig
}

func GetCliVersion() string {
	filePath := "cocli_version.txt"
	_, err := os.Stat(filePath)
	if err == nil {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		content, err := io.ReadAll(file)
		if err != nil {
			panic(err)
		}

		if string(content) != cliVersion {
			panic("cocli version mismatch - make sure constant cliVersion and the contents of cocli_version.txt match!")
		}
	}

	return cliVersion
}

func GetCoherenceDomain() string {
	domain, domain_exists := os.LookupEnv("COHERENCE_DOMAIN")
	if domain_exists {
		return domain
	}

	return GetCliConfig().CoherenceDomain
}

func CoherenceApiRequest(method string, url string, body io.Reader) (*http.Response, error) {
	return AuthenticatedRequest(
		method,
		url,
		body,
		os.Getenv("COHERENCE_ACCESS_TOKEN"),
	)
}

func AuthenticatedRequest(method string, url string, body io.Reader, bearer_token string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest(
		method,
		url,
		body,
	)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", bearer_token))

	res, err := client.Do(req)

	if res != nil && res.StatusCode == 401 {
		fmt.Println("Unauthorized... COHERENCE_ACCESS_TOKEN missing or expired/invalid. Update the value of the COHERENCE_ACCESS_TOKEN and try again.")
	}

	if res.StatusCode != http.StatusOK {
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(res.StatusCode)
		panic(FormatJSONOutput(bodyBytes))
	}

	return res, err
}

func FormatFilename(filename string) string {
	if filename[:2] == "~/" {
		usr, err := user.Current()
		if err != nil {
			panic(err)
		}
		filename = filepath.Join(usr.HomeDir, filename[2:])
	}

	return filename
}

func FormatJSONOutput(bodyBytes []byte) string {
	var formattedRespBody bytes.Buffer
	err := json.Indent(&formattedRespBody, bodyBytes, "", "    ")
	if err != nil {
		panic(err)
	}

	return formattedRespBody.String()
}
