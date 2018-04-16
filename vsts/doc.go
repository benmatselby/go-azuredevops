/*
go-vsts is a Go client library for accessing the Visual Studio Team Services API.
Installation
	$ go get github.com/benmatselby/go-vsts/vsts
	or
	$ dep ensure -add github.com/benmatselby/go-vsts/vsts
Usage
Interaction with the VSTS API is done through a Client instance.
	import "github.com/benmatselby/go-vsts/vsts
	v := vsts.NewClient(account, project, token)
Services
The client has services that you can use to access resources from the API:
	iterations, error := v.Iterations.List(team)
	if error != nil {
		fmt.Println(error)
	}
*/
package vsts
