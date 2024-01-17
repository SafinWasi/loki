# loki

![Tests](https://github.com/SafinWasi/loki/actions/workflows/test.yaml/badge.svg) ![Builds](https://github.com/SafinWasi/loki/actions/workflows/build.yaml/badge.svg)

A simple OpenID Connect Requesting Party (RP) to authenticate with an OpenID Provider, written fully in go.

```
Usage:
  loki [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  run         Runs the Loki server

Flags:
      --debug   Debug mode
  -h, --help    help for loki

Use "loki [command] --help" for more information about a command.
```

## Requirements

Requires go 1.20

## Building

Run `go build` to build an executable. Binary releases will be available soon.

## Testing

```
go test ./...
```
