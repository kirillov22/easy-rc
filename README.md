# Mouse Server

Soon to be a single monorepo that contains the client, server and shared protobuf model for the two apps.

## Prerequisites
- Golang version v1.25 - https://go.dev/dl/
  - Using a version manager such as goenv is recommended - https://github.com/go-nv/goenv
- Protobuf compiler - https://protobuf.dev/installation/
- `protoc-gen-go` - Go code generator plugin for protobuf
- NodeJS v24.12.0 - https://nodejs.org/en/download
  - Using a version manager such as nvm is recommended - https://github.com/nvm-sh/nvm

### Setting up Go with goenv

1. Install goenv: https://github.com/go-nv/goenv#installation
2. Install and activate the required Go version:
   ```bash
   goenv install 1.25.5
   goenv global 1.25.5   # or `goenv local 1.25.5` inside this repo
   ```
3. Ensure `$(go env GOPATH)/bin` is on your PATH. Add this to your shell profile (`.zshrc`, `.bashrc`, etc.):
   ```bash
   export PATH="$(go env GOPATH)/bin:$PATH"
   ```
4. Install the protobuf Go plugin:
   ```bash
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   ```


## Server
- Go application with an http connection which will serve the client js code
- Web socket connection for the js client to talk to in order to issue commands for the server to read and execute on the machine
- The server will provide all instructions to the client how it can find the specific instance of the device (which ip address and port to connect to)

## Client
- Basic JS application that acts as a remote from your mobile device to your server

## Common
- Shared protobuf model for the client and server to communicate with one another

# Scripts

Build scripts are in the `scripts/` directory. They can be run from any working directory.

| Script | Description |
|--------|-------------|
| `./scripts/proto.sh` | Regenerate Go code from the protobuf schema |
| `./scripts/dev.sh` | Build the server with debug logging enabled |
| `./scripts/build.sh` | Production build: generates proto, runs vet + tests, then builds the binary |

All scripts output the server binary to `server/bin`.

# Running

## Server
- Development: `./scripts/dev.sh && ./server/bin/easy-rc`
- Production: `./scripts/build.sh && ./server/bin/easy-rc`

## TODO list
- [x] Migrate from the custom-built protocol to the protobuf schema
- [x] Migrate the code (merge client) into a monorepo to easier share the protobuf models
- [x] Create QR code/ simple way for client to connect - https://github.com/caseymrm/menuet
- [x] Refactor server to have the websocket be decoupled from the command processing so it can be tested
- [ ] Create scripts that build & run the applications all at once
- [ ] Get CI/CD pipeline working in github actions
- [ ] Implement double clicking
- [ ] Implement clicking on the touchpad instead of relying on the buttons
- [ ] Fix reconnection bug after refreshing the browser it doesn't reconnect straight away
- [ ] Rewrite client to some frontend framework. Getting quite tiresome manually managing the state even though it is small 
