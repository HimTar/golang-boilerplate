If you get any error saying "air: File or directory not found", follow below steps

1. Export the GOPATH/bin directory.

```
export PATH=$PATH:$GOPATH/bin
```

Note : Might be possible the GOPATH is not set, in that case set to `/Users/:username/go`

If you are not getting any intellisence on local or imported packages, follow below steps

1. Press `cmd + shift + P`.
2. Search for `go tools`
3. Install/Update them
