# SteerMesh Go SDK

Go client for the SteerMesh Cloud API. Use for pipelines, automation, or the steer CLI.

## Install

```bash
go get github.com/SteerMesh/go-sdk
```

## Usage

```go
import "github.com/SteerMesh/go-sdk/client"

c := client.New("https://api.steermesh.dev", "your-api-key")
packs, err := c.ListPacks()
bundle, err := c.GetBundle("bundle-id")
err = c.DownloadBundle("bundle-id", "./output")
```

## License

See repository license.
