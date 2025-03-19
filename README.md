# recallai-go

[![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/harrison-peng/recallai-go?label=go%20module)](https://github.com/harrison-peng/recallai-go/tags)
[![Go Reference](https://pkg.go.dev/badge/github.com/harrison-peng/recallai-go.svg)](https://pkg.go.dev/github.com/harrison-peng/recallai-go)
[![Test](https://github.com/harrison-peng/recallai-go/actions/workflows/test.yml/badge.svg)](https://github.com/harrison-peng/recallai-go/actions/workflows/test.yml)

`recallai-go` is a Go client library for interacting with the [Recall.ai API](https://docs.recall.ai/).

## Supported APIs

This library is continually updated to support all endpoints of the Recall.ai API.

## Installation

```bash
go get github.com/harrison-peng/recallai-go
```

## Usage

First, please follow the [Getting Started Guide](https://docs.recall.ai/docs/getting-started) to obtain an integration token.

### Initialization

Import this library and initialize the API client using the obtained integration token.

```go
import "github.com/harrison-peng/recallai-go"

client := recallaigo.NewClient("your_integration_token")
```

### Calling the API

You can use the methods of the initialized client to call the Recall.ai API. Here is an example of how to list all bots:

```go
page, err := client.Bot.List(context.Background(), nil)
if err != nil {
    // Handle the error
}
```
