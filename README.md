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
