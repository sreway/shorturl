# cmd/shortener

## Build info

Setup the -ldflags option for build info

Build:
```sh
go build -v -ldflags="-X 'main.buildVersion=0.1.0' \
 -X 'main.buildDate=$(date)' \
 -X 'main.buildCommit=test'" -o ./shortener \
 cmd/shortener/main.go
```
Run:
```sh
go run -v -ldflags="-X 'main.buildVersion=0.1.0' \
 -X 'main.buildDate=$(date)' \
 -X 'main.buildCommit=test'" \
cmd/shortener/main.go
```

