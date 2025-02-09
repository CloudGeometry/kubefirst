/*
Copyright (C) 2021-2024, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.
*/
package azure

import (
	"fmt"
	"os"

	internalssh "github.com/konstructio/kubefirst-api/pkg/ssh"
	"github.com/rs/zerolog/log"
)

// Environment variables required for authentication. This should be a
// service principal - the Terraform provider docs detail how to create
// one
// @link https://registry.terraform.io/providers/hashicorp/azurerm/latest/docs/guides/service_principal_client_secret.html
var envvarSecrets = []string{
	"ARM_CLIENT_ID",
	"ARM_CLIENT_SECRET",
	"ARM_TENANT_ID",
	"ARM_SUBSCRIPTION_ID",
}

func ValidateProvidedFlags(gitProvider string) error {
	for _, env := range envvarSecrets {
		if os.Getenv(env) == "" {
			return fmt.Errorf("your %s is not set - please set and re-run your last command", env)
		}
	}

	switch gitProvider {
	case "github":
		key, err := internalssh.GetHostKey("github.com")
		if err != nil {
			return fmt.Errorf("known_hosts file does not exist - please run `ssh-keyscan github.com >> ~/.ssh/known_hosts` to remedy")
		} else {
			log.Info().Msgf("%s %s\n", "github.com", key.Type())
		}
	case "gitlab":
		key, err := internalssh.GetHostKey("gitlab.com")
		if err != nil {
			return fmt.Errorf("known_hosts file does not exist - please run `ssh-keyscan gitlab.com >> ~/.ssh/known_hosts` to remedy")
		} else {
			log.Info().Msgf("%s %s\n", "gitlab.com", key.Type())
		}
	}

	return nil
}
