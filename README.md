# go-azuredevops

[![GoDoc](https://godoc.org/github.com/benmatselby/go-azuredevops/azuredevops?status.svg)](https://godoc.org/github.com/benmatselby/go-azuredevops/azuredevops)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=go-azuredevops&metric=alert_status)](https://sonarcloud.io/dashboard?id=go-azuredevops)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/go-azuredevops?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/go-azuredevops)

`go-azuredevops` is a Go client library for accessing the [Azure DevOps API](https://docs.microsoft.com/en-gb/rest/api/vsts/). This is very much work in progress, so at the moment supports a small subset of the API.

## Services

There is partial implementation for the following services

- Boards
- Builds
- Favourites
- Iterations
- Pull Requests
- Work Items

## Usage

```go
import "github.com/benmatselby/go-azuredevops/azuredevops
```

Construct a new Azure DevOps Client

```go
v := azuredevops.NewClient(account, project, token)
```

Get a list of iterations

```go
iterations, error := v.Iterations.List(team)
if error != nil {
    fmt.Println(error)
}

for index := 0; index < len(iterations); index++ {
    fmt.Println(iterations[index].Name)
}
```
