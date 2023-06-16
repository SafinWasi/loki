# Loki

A simple OpenID Connect Requesting Party (RP) to authenticate with an OpenID Provider, written fully in go.

```
Usage:
  loki [command] [flags]
  loki [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  delete      Deletes a configuration by alias
  help        Help about any command
  list        List aliases configured
  register    Registers a new OpenID Client
  setup       Setup details for OIDC

Flags:
  -a, --alias string   Flow to be used for authentication
      --disable-ssl    Disables SSL
  -f, --flow string    Flow to be used for authentication
  -h, --help           help for loki

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
