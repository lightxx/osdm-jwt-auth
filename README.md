# osdm-jwt-auth

`osdm-jwt-auth` creates a JWT client assertion, exchanges it against an OpenID Connect token endpoint, and prints the returned `access_token` to stdout.

The tool is intended for scriptable use:

```sh
ACCESS_TOKEN="$(./osdm-jwt-auth ...)"
```

If you pass `-verbose`, the tool writes the intermediate client assertion and decoded assertion details to stderr while keeping stdout clean for the final access token.

You can either download the binaries from the latest release package (automatically created by a github runner), or, like you should with everything crypto related, build it yourself. See instructions below. 

## Requirements

- Go `1.26.1` or compatible
- An RSA private key in PEM format
- Network access to the target token endpoint

## Build

The project is a single Go module and can be built with `go build`.

### Build on macOS

```sh
go build -o osdm-jwt-auth .
```

### Build on Linux

```sh
go build -o osdm-jwt-auth .
```

### Build on Windows

In PowerShell:

```powershell
go build -o osdm-jwt-auth.exe .
```

## Cross-compile from any OS

### Build a Linux binary

```sh
GOOS=linux GOARCH=amd64 go build -o osdm-jwt-auth-linux-amd64 .
```

### Build a macOS binary

```sh
GOOS=darwin GOARCH=arm64 go build -o osdm-jwt-auth-darwin-arm64 .
```

### Build a Windows binary

```sh
GOOS=windows GOARCH=amd64 go build -o osdm-jwt-auth-windows-amd64.exe .
```

## Generate an RSA key pair

The tool expects a PEM-encoded RSA private key. The examples below generate a 2048-bit key pair.

### Linux

Using OpenSSL:

```sh
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out osdm_private.pem
openssl rsa -pubout -in osdm_private.pem -out osdm_public.pem
```

### macOS

Using OpenSSL:

```sh
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out osdm_private.pem
openssl rsa -pubout -in osdm_private.pem -out osdm_public.pem
```

### Windows

In PowerShell with OpenSSL on `PATH`:

```powershell
openssl genpkey -algorithm RSA -pkeyopt rsa_keygen_bits:2048 -out osdm_private.pem
openssl rsa -pubout -in osdm_private.pem -out osdm_public.pem
```

If OpenSSL is not installed, one common option is to install it with `winget` first:

```powershell
winget install ShiningLight.OpenSSL
```

Then reopen PowerShell and run the key generation commands above.

## Usage

The tool:

1. Creates a signed JWT client assertion.
2. Sends it to the token endpoint using the client credentials flow.
3. Prints the returned `access_token` to stdout.

### Flags

- `-key`: path to the RSA private key PEM file
- `-alg`: signing algorithm, `RS256` or `PS256`
- `-iss`: issuer / client ID used in the assertion
- `-sub`: subject / client ID used in the assertion
- `-aud`: token endpoint URL used as both JWT audience and token exchange target
- `-scope`: scope sent in the assertion and token exchange request
- `-lifetime`: assertion lifetime in seconds
- `-nbf`: include `nbf` in the assertion
- `-iat`: include `iat` in the assertion
- `-grace`: seconds subtracted from `now` when `-nbf` is enabled
- `-verbose`: print the generated client assertion and decoded assertion details to stderr

## Example: Get an access token

- client ID: `osdm_test09`
- token endpoint: `https://apim-yvorp.om.tsint.at/auth/realms/osdm/protocol/openid-connect/token`
- private key: `./osdm_yvorp_tom_private.pem`
- lifetime: `300`

### macOS or Linux

```sh
ACCESS_TOKEN="$(./osdm-jwt-auth \
  -key ./osdm_yvorp_tom_private.pem \
  -iss osdm_test09 \
  -sub osdm_test09 \
  -lifetime 300 \
  -scope openid \
  -aud https://apim-yvorp.om.tsint.at/auth/realms/osdm/protocol/openid-connect/token)"
```

Print the token:

```sh
printf '%s\n' "$ACCESS_TOKEN"
```

Verbose mode:

```sh
ACCESS_TOKEN="$(./osdm-jwt-auth \
  -key ./osdm_yvorp_tom_private.pem \
  -iss osdm_test09 \
  -sub osdm_test09 \
  -lifetime 300 \
  -scope openid \
  -aud https://apim-yvorp.om.tsint.at/auth/realms/osdm/protocol/openid-connect/token \
  -verbose)"
```

### Windows PowerShell

```powershell
$ACCESS_TOKEN = .\osdm-jwt-auth.exe `
  -key .\osdm_yvorp_tom_private.pem `
  -iss osdm_test09 `
  -sub osdm_test09 `
  -lifetime 300 `
  -scope openid `
  -aud https://apim-yvorp.om.tsint.at/auth/realms/osdm/protocol/openid-connect/token
```

Print the token:

```powershell
$ACCESS_TOKEN
```

Verbose mode:

```powershell
$ACCESS_TOKEN = .\osdm-jwt-auth.exe `
  -key .\osdm_yvorp_tom_private.pem `
  -iss osdm_test09 `
  -sub osdm_test09 `
  -lifetime 300 `
  -scope openid `
  -aud https://apim-yvorp.om.tsint.at/auth/realms/osdm/protocol/openid-connect/token `
  -verbose
```

## Output behavior

- Default mode: stdout contains only the exchanged `access_token`
- Verbose mode: stdout still contains only the `access_token`; stderr additionally shows the generated client assertion plus assertion header and payload details
