# dedicated-go
Go SDK for Selectel Dedicated Servers

## Getting started

### Installation

You can install needed `dedicated-go` packages via `go get` command:

```bash
go get github.com/selectel/dedicated-go
```

### Authentication

To work with the Selectel Dedicated Servers API you first need to:

* Create a Selectel account: [registration page](https://my.selectel.ru/registration).
* Create a project in Selectel Cloud Platform [projects](https://my.selectel.ru/vpc/projects).
* Retrieve a token for your project via API or [go-selvpcclient](https://github.com/selectel/go-selvpcclient).

### Endpoints

You can find available endpoints [here](https://docs.selectel.ru/en/api/urls/).

### Usage example

```go
package main

import (
	"context"
	"fmt"
	"log"

	dedicated "github.com/selectel/dedicated-go/pkg/v2"
)

func main() {
	// Auth token.
	token := "gAAAAABeVNzu-..."

	// Dedicated servers endpoint to work with.
	endpoint := "https://api.selectel.ru/servers/v2"

	// Create the client.
	client := dedicated.NewClientV2(token, endpoint)

	// Get the os for location.
	q := &dedicated.OperatingSystemsQuery{LocationID: "some-location-uuid"}
	operatingSystems, _, err := client.OperatingSystems(context.Background(), q)
	if err != nil {
		log.Fatal(err)
	}

	// Print the os.
	for idx, os := range operatingSystems {
		fmt.Printf("OS %d: %+v", idx, os)
	}
}
```