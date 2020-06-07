# Jobatator CLI Client

This is a little tool to allow human interaction with a jobatator instance throught the command line interface.

![GIF Demo](https://media.giphy.com/media/THgR7WPa9gYmRQLSC9/giphy.gif)

## Installation

### Using docker

`docker run jobatator/cli your_username:your_password@your_host:8962`

If you want to access a jobatator instance that is running on your host machine it's a little bit more complicated, if you are on linux the `--network="host"` flag may help but it's not working for me. Check these stackoverflow topics: [How to access host port from docker container](https://stackoverflow.com/questions/31324981/how-to-access-host-port-from-docker-container#43541732); [From inside of a Docker container, how do I connect to the localhost of the machine?](https://stackoverflow.com/questions/24319662/from-inside-of-a-docker-container-how-do-i-connect-to-the-localhost-of-the-mach#24326540)

### Build from source

- `git clone https://github.com/jobatator/cli.git jobatator-cli`
- `cd jobatator-cli`
- `go build -o jobatator-cli main.go`
- Your binary is ready to be used at the following location: `jobatator-cli/jobatator-cli`

## Usage

Usage:	jobatator-cli [OPTIONS] URI

Will connect to the jobatator instance

The URI is use to specify the host, port, username, password and group of the session:

`[[username][:password]@]host[:port][/group]`

Flags:

- `-r` or `--raw` : Will disable the automatic JSON formatting and will instead show the raw ugly JSON

## Features

- Automatic connection using the `AUTH` command
- Automatic group selection using the `USE_GROUP` command
- Command autocompletion
- Command history (Up & Down arrow keys)
- JSON formatting
- A more human way of communication with a jobatator instance

Side note: In fact, you could use the [netcat](https://en.wikipedia.org/wiki/Netcat) command it will totaly do the work but this cli provide some extra feature to feel more confortable.

## Credits

- For auto completion, command history this cli is using the [c-bata/go-prompt](https://github.com/c-bata/go-prompt) lib
- For JSON formatting, this use the [TylerBrock/colorjson](github.com/TylerBrock/colorjson) lib

