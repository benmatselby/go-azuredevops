go-vsts
=======

[![Build Status](https://travis-ci.org/benmatselby/go-vsts.png?branch=master)](https://travis-ci.org/benmatselby/go-vsts)
[![Go Report Card](https://goreportcard.com/badge/github.com/benmatselby/go-vsts?style=flat-square)](https://goreportcard.com/report/github.com/benmatselby/go-vsts)

go-vsts is a Go client library for accessing the [Visual Studio Team Services API](https://docs.microsoft.com/en-gb/rest/api/vsts/). This is very much work in progress, so at the moment supports a small subset of the API.

## Services

There is partial implementation for the following services

* Boards
* Builds
* Favourites
* Iterations
* Pull Requests
* Work Items

## Usage

```go
import "github.com/benmatselby/go-vsts/vsts
```

Construct a new VSTS Client

```go
v := vsts.NewClient(account, project, token)
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
