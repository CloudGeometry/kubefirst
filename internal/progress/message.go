/*
Copyright (C) 2021-2023, Kubefirst

This program is licensed under MIT.
See the LICENSE file for more details.

Emojis definition https://github.com/yuin/goldmark-emoji/blob/master/definition/github.go
Color definition https://www.ditig.com/256-colors-cheat-sheet
*/
package progress

import (
	"fmt"
	"log"
	"strconv"

	"github.com/charmbracelet/glamour"
	"github.com/konstructio/kubefirst-api/pkg/types"
	"github.com/spf13/viper"
)

func RenderMessage(message string) string {
	r, _ := glamour.NewTermRenderer(
		glamour.WithStyles(StyleConfig),
		glamour.WithEmoji(),
	)

	out, err := r.Render(message)
	if err != nil {
		log.Println(err.Error())
		return err.Error()
	}
	return out
}

func createStep(message string) addStep {
	out := RenderMessage(message)

	return addStep{
		message: out,
	}
}

func createErrorLog(message string) errorMsg {
	out := RenderMessage(fmt.Sprintf("##### :no_entry_sign: Error: %s", message))

	return errorMsg{
		message: out,
	}
}

// Public Progress Functions
func DisplayLogHints(estimatedTime int) {
	logFile := viper.GetString("k1-paths.log-file")
	cloudProvider := viper.GetString("kubefirst.cloud-provider")

	documentationLink := "https://kubefirst.konstruct.io/docs/"
	if cloudProvider != "" {
		documentationLink += cloudProvider + `/quick-start/install/cli`
	}

	header := `
##
# Welcome to Kubefirst

### :bulb: To view verbose logs run below command in new terminal:
` + fmt.Sprintf("##### **tail -f -n +1 %s**", logFile) + `
### :blue_book: Documentation: ` + documentationLink + `

### :alarm_clock: Estimated time:` + fmt.Sprintf("`%s minutes` \n\n", strconv.Itoa(estimatedTime))

	headerMessage := RenderMessage(header)

	Progress.Send(headerMsg{
		message: headerMessage,
	})
}

//nolint:revive // will be fixed in the future
func DisplaySuccessMessage(cluster types.Cluster) string {
	cloudCliKubeconfig := ""

	gitProviderLabel := "GitHub"
	if cluster.GitProvider == "gitlab" {
		gitProviderLabel = "GitLab"
	}

	switch cluster.CloudProvider {
	case "aws":
		cloudCliKubeconfig = fmt.Sprintf("aws eks update-kubeconfig --name %q --region %q", cluster.ClusterName, cluster.CloudRegion)
	case "azure":
		cloudCliKubeconfig = fmt.Sprintf("az aks get-credentials --resource-group %q --name %q", cluster.ClusterName, cluster.ClusterName)
	case "civo":
		cloudCliKubeconfig = fmt.Sprintf("civo kubernetes config %q --save", cluster.ClusterName)
	case "digitalocean":
		cloudCliKubeconfig = "doctl kubernetes cluster kubeconfig save " + cluster.ClusterName
	case "google":
		cloudCliKubeconfig = fmt.Sprintf("gcloud container clusters get-credentials %q --region=%q", cluster.ClusterName, cluster.CloudRegion)
	case "vultr":
		cloudCliKubeconfig = fmt.Sprintf("vultr-cli kubernetes config %q", cluster.ClusterName)
	case "k3s":
		cloudCliKubeconfig = "use the kubeconfig file outputted from terraform to access the cluster"
	}

	var fullDomainName string
	if cluster.SubdomainName != "" {
		fullDomainName = fmt.Sprintf("%s.%s", cluster.SubdomainName, cluster.DomainName)
	} else {
		fullDomainName = cluster.DomainName
	}

	success := `
##
#### :tada: Success` + "`Cluster " + cluster.ClusterName + " is now up and running`" + `

# Cluster ` + cluster.ClusterName + ` details:

### :bulb: To retrieve root credentials for your Kubefirst platform run:
##### kubefirst ` + cluster.CloudProvider + ` root-credentials

## ` + fmt.Sprintf("`%s `", gitProviderLabel) + `
### Git Owner   ` + fmt.Sprintf("`%s`", cluster.GitAuth.Owner) + `
### Repos       ` + fmt.Sprintf("`https://%s.com/%s/gitops` \n\n", cluster.GitProvider, cluster.GitAuth.Owner) +
		fmt.Sprintf("`            https://%s.com/%s/metaphor`", cluster.GitProvider, cluster.GitAuth.Owner) + `
## Kubefirst Console
### URL         ` + fmt.Sprintf("`https://kubefirst.%s`", fullDomainName) + `
## Argo CD
### URL         ` + fmt.Sprintf("`https://argocd.%s`", fullDomainName) + `
## Vault
### URL         ` + fmt.Sprintf("`https://vault.%s`", fullDomainName) + `


### :bulb: Quick start examples:

### To connect to your new Kubernetes cluster run:
##### ` + cloudCliKubeconfig + `

### To view all cluster pods run:
##### kubectl get pods -A
`

	return success
}

func AddStep(message string) {
	renderedMessage := createStep(fmt.Sprintf("%s %s", ":dizzy:", message))
	Progress.Send(renderedMessage)
}

func CompleteStep(message string) {
	Progress.Send(completeStep{
		message: message,
	})
}

func Success(success string) {
	successMessage := RenderMessage(success)

	Progress.Send(
		successMsg{
			message: successMessage,
		})
}

func Error(message string) {
	renderedMessage := createErrorLog(message)
	Progress.Send(renderedMessage)
}

func StartProvisioning(clusterName string) {
	provisioningMessage := startProvision{
		clusterName: clusterName,
	}

	Progress.Send(provisioningMessage)
}
