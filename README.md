# Trippy

A Slot Machine which run with the Atkins Diet Machine as the engine and has support for other types of Slot machines.


## Usage

To use Trippy, Go 1.10 or higher needs to be installed. [Install Go](https://golang.org/doc/install)

1. Clone the repository into ~/go/src

2. Ensure that either GOPATH is empty or GOPATH=~/go

3. Run `cd ~/go/src/trippy && go install ./...`

4. A secret API key is used to encrypt JWT token. The key is provided to the server in a file.
   Start the server with the API key filepath in env.

   `echo "abc" >> keyfile && export TRIPPY_API_KEY_PATH=./keyfile && trippy`
