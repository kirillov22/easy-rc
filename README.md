# Mouse Server

Soon to be a single monorepo that contains the client, server and shared protobuf model for the two apps.

## Prerequisites
- Golang version v1.25 - https://go.dev/dl/
  - Using a version manager such as goenv is recommended - https://github.com/go-nv/goenv
- Protobuf compiler - https://protobuf.dev/installation/
- NodeJS v24.12.0 - https://nodejs.org/en/download
  - Using a version manager such as nvm is recommended - https://github.com/nvm-sh/nvm


## Server
- Go application with an http connection which will serve the client js code
- Web socket connection for the js client to talk to in order to issue commands for the server to read and execute on the machine
- The server will provide all instructions to the client how it can find the specific instance of the device (which ip address and port to connect to)

## Client
- Basic JS application that acts as a remote from your mobile device to your server

## Common
- Shared protobuf model for the client and server to communicate with one another

# Building
- Build the protobuf `protoc -I=/<PATH_TO_REPO>/easy-rc --go_out=/<PATH_TO_REPO>/easy-rc common/messages.proto`

# Running

## TODO list
- [ ] Migrate from the custom-built protocol to the protobuf schema
- [ ] Migrate the code (merge client) into a monorepo to easier share the protobuf models
- [ ] Create QR code/ simple way for client to connect - https://github.com/caseymrm/menuet
- [ ] Refactor server to have the websocket be decoupled from the command processing so it can be tested
- [ ] Create scripts that build & run the applications all at once
- [ ] Tests xD
