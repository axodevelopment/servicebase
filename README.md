# NOT MEANT TO BE USED

This is just a place I am testing some mongo modules I am putting together. I use this to rapidly onboard go based apps for kubernetes use. Very alpha very incomplete but i'll be adding more to t and cleaning things up. Not meant for prod use just demo boilerplate stuff.

## Table of Contents

- [General Usage](#general_usage)

## General Usage

```go
import (
  serviceBase "github.com/axodevelopment/servicebase"
)


func main() {

  var svc *serviceBase.Service

  svc, _ = serviceBase.New("AirportApp", serviceBase.WithPort(9091), serviceBase.WithHealthProbe(true))


	/*
    Do whatever you want
  */

  svc.AppHealthz = true
  svc.AppReadyz = true

	//start the backend
	go func(s *serviceBase.Service) {
		serviceBase.Start(s)
	}(svc)

	<-svc.ExitAppChan
}
```

Goal here is this will wrap up GIN and health / ready probes etc so you can just to get to making routes etc.

```go

func createAirportRoutes(svc *serviceBase.Service) {
	svc.GinEngine.GET("/route", func(ctx *gin.Context) {

		/// do stuff
    ctx.JSON(http.StatusOK, data)

	})
}

```
