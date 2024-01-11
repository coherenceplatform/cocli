package cocli

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	"github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

type UseExisting struct {
	Image string `json:"image"`
	Mode  string `json:"mode"`
	Tag   string `json:"tag"`
}

type ProdImgConfigOverrideItem struct {
	ServiceName string `json:"service_name"`
	Prod        struct {
		UseExisting UseExisting `json:"use_existing"`
	} `json:"prod"`
}

type DevImgConfigOverrideItem struct {
	ServiceName string `json:"service_name"`
	Dev         struct {
		UseExisting UseExisting `json:"use_existing"`
	} `json:"dev"`
}

type ImageOverridePayload struct {
	ConfigOverrideItems []interface{} `json:"config_override_items"`
}

var envImageOverrideEnvId int
var envImageOverrideServiceName string

var envImageOverrideImage string
var envImageOverrideTag string
var envImageOverrideMode string
var envImageOverrideDev bool

var environmentImageOverrideCmd = &cobra.Command{
	Use:   "set_service_image",
	Short: "Override service image configuration",
	Long: "Overrides the image configuration for services " +
		"that use an existing image (not built by coherence)",
	Run: func(cmd *cobra.Command, args []string) {
		envImageOverridePayloadData := ImageOverridePayload{
			ConfigOverrideItems: []interface{}{
				ProdImgConfigOverrideItem{
					ServiceName: envImageOverrideServiceName,
					Prod: struct {
						UseExisting UseExisting `json:"use_existing"`
					}{
						UseExisting: UseExisting{
							Image: envImageOverrideImage,
							Mode:  envImageOverrideMode,
							Tag:   envImageOverrideTag,
						},
					},
				},
			},
		}

		if envImageOverrideDev == true {
			envImageOverridePayloadData = ImageOverridePayload{
				ConfigOverrideItems: []interface{}{
					DevImgConfigOverrideItem{
						ServiceName: envImageOverrideServiceName,
						Dev: struct {
							UseExisting UseExisting `json:"use_existing"`
						}{
							UseExisting: UseExisting{
								Image: envImageOverrideImage,
								Mode:  envImageOverrideMode,
								Tag:   envImageOverrideTag,
							},
						},
					},
				},
			}
		}

		payloadBytes, err := json.Marshal(envImageOverridePayloadData)
		if err != nil {
			panic(err)
		}
		payload := bytes.NewBuffer(payloadBytes)

		envConfigOverrideUrl := fmt.Sprintf(
			"https://%s%s/environments/%s/update_config_overrides",
			cocli.GetCoherenceDomain(),
			cocli.GetCoherenceApiPrefix(),
			fmt.Sprint(envImageOverrideEnvId),
		)
		res, err := cocli.CoherenceApiRequest(
			"PUT",
			envConfigOverrideUrl,
			payload,
		)
		if err != nil {
			panic(err)
		}
		defer res.Body.Close()
		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		fmt.Println(cocli.FormatJSONOutput(bodyBytes))
	},
}

func init() {
	environmentImageOverrideCmd.Flags().IntVarP(&envImageOverrideEnvId, "environment_id", "e", 0, "Environment ID (required)")
	environmentImageOverrideCmd.MarkFlagRequired("environment_id")

	environmentImageOverrideCmd.Flags().StringVarP(&envImageOverrideServiceName, "service", "s", "", "Service name (required)")
	environmentImageOverrideCmd.MarkFlagRequired("service")

	environmentImageOverrideCmd.Flags().StringVarP(&envImageOverrideImage, "image", "i", "", "Image name without tag (optional)")
	environmentImageOverrideCmd.Flags().StringVarP(&envImageOverrideTag, "tag", "t", "", "Image tag (optional)")
	environmentImageOverrideCmd.Flags().StringVarP(&envImageOverrideMode, "mode", "m", "", "Image mode (optional)")

	environmentImageOverrideCmd.Flags().BoolVarP(
		&envImageOverrideDev,
		"dev", "d", false,
		"Only apply settings to the dev image config - this is the image used in workspaces (optional, default: false)",
	)
}
