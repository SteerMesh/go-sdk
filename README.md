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
if err != nil {
  // handle error
}

bundle, err := c.GetBundle("bundle-id")
if err != nil {
  // handle error
}

err = c.DownloadBundle("bundle-id", "./output")
if err != nil {
  // handle error
}
```

## Compatibility
The SDK tracks the public SteerMesh Cloud API. Match SDK versions with the Cloud API version used by your org.

## License

See repository license.
