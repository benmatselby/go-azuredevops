go-vsts
=======

[![Build Status](https://travis-ci.org/benmatselby/go-vsts.png?branch=master)](https://travis-ci.org/benmatselby/go-vsts)

go-vsts is a Go client library for accessing the [VSTS API](https://docs.microsoft.com/en-gb/rest/api/vsts/). This is very much work in progress, so at the moment supports a small subset of the API.


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
