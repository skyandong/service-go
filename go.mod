module github.com/skyandong/service-go

go 1.15

replace github.com/skyandong/autils => ../autils

require (
	github.com/gin-gonic/gin v1.7.7
	github.com/skyandong/autils v0.0.1
	go.uber.org/zap v1.18.1
)
