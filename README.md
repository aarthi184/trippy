# Trippy

A Slot Machine which run with the Atkins Diet Machine as the engine and has support for other types of Slot machines.


## Usage

A secret API key is used to encrypt JWT token. The key is provided to the server in a file.
Start the server with the API key filepath in env.

`echo "abc" >> keyfile && export TRIPPY_API_KEY_PATH=./keyfile && trippy`
