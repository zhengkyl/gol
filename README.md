# gtg

git terminal go

Add to `~/.ssh/config` or manually clear out `localhost` entries `~/.ssh/known_hosts`. Otherwise, it may cause issues when connecting to local servers in development with different server keys.

```
Host localhost
    UserKnownHostsFile /dev/null
```
