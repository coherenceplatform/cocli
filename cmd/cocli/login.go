package cocli

import (
	"context"
	"fmt"

	cocli "github.com/coherenceplatform/cocli/pkg/cocli"
	"github.com/spf13/cobra"
)

var loginCmd = &cobra.Command{
	Use:   "login",
	Short: "Authenticate with the Coherence API",
	Run: func(cmd *cobra.Command, args []string) {
		// Check if creds configured via env var
		if cocli.IsRefreshTokenVarSet() {
			fmt.Println("\nAlready logged in... to logout unset COCLI_REFRESH_TOKEN")
			return
		}

		// Check if there are already creds present (maybe in some config file)
		token := cocli.GetTokenFromCredsFile()
		if token != nil {
			fmt.Println("\nAlready logged in... to logout use `cocli logout`")
			return
		}
		// If not, start device code flow
		// Poll for device code flow completion or timeout
		// If authenticated store access & refresh token in some .config file
		config := cocli.GetOauthConfig()
		ctx := context.Background()
		response, err := config.DeviceAuth(ctx)
		if err != nil {
			panic(err)
		}
		fmt.Printf("Please login at %s\n", response.VerificationURIComplete)
		token, err = config.DeviceAccessToken(ctx, response)
		if err != nil {
			panic(err)
		}
		tokenWithIDToken := &cocli.TokenWithIdToken{
			Token:   token,
			IDToken: token.Extra("id_token").(string),
		}

		cocli.WriteTokenFile(tokenWithIDToken)
	},
}
