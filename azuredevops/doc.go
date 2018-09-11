/*
Package azuredevops is a Go client library for accessing the Azure DevOps API.
Installation
	$ go get github.com/benmatselby/go-azuredevops/azuredevops
	or
	$ dep ensure -add github.com/benmatselby/go-azuredevops/azuredevops
Usage
Interaction with the Azure DevOps API is done through a Client instance.
	import "github.com/benmatselby/go-azuredevops/azuredevops
	v := azuredevops.NewClient(account, project, token)
Services
The client has services that you can use to access resources from the API:
	iterations, error := v.Iterations.List(team)
	if error != nil {
		fmt.Println(error)
	}
*/
package azuredevops
