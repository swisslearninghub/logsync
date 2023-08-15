module github.com/swisslearninghub/logsync

go 1.18

replace (
	github.com/swisslearninghub/logsync/api => ./api
	github.com/swisslearninghub/logsync/cefsyslog => ./cefsyslog
	github.com/swisslearninghub/logsync/commands => ./commands
	github.com/swisslearninghub/logsync/config => ./config
)

require (
	github.com/go-playground/validator/v10 v10.14.0
	github.com/urfave/cli/v2 v2.25.4
	golang.org/x/oauth2 v0.8.0
)

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.2 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/leodido/go-urn v1.2.4 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/xrash/smetrics v0.0.0-20201216005158-039620a65673 // indirect
	golang.org/x/crypto v0.7.0 // indirect
	golang.org/x/net v0.10.0 // indirect
	golang.org/x/sys v0.8.0 // indirect
	golang.org/x/text v0.9.0 // indirect
	google.golang.org/appengine v1.6.7 // indirect
	google.golang.org/protobuf v1.28.0 // indirect
)
