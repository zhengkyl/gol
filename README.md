# gol

game of life using [wish](https://github.com/charmbracelet/wish) and [bubbletea](https://github.com/charmbracelet/bubbletea)

## Running locally

```sh
go run main.go

# in a separate terminal
ssh -p2345 localhost
```

## IMPORTANT!

Add to `~/.ssh/config` or manually clear out `localhost` entries `~/.ssh/known_hosts`. Otherwise, it may cause issues b/c the server's key signature changes each time it restarts.

```
Host localhost
    UserKnownHostsFile /dev/null
```

### Building

If you see a `no such file or directory` error when running the container, see try this.

https://stackoverflow.com/questions/36279253/go-compiled-binary-wont-run-in-an-alpine-docker-container-on-ubuntu-host

```
CGO_ENABLED=0 go build
```
