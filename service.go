package servicebase

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

type Service struct {
	//public
	Name        string          `json:"Name"`
	Port        int             `json:"Port"`
	HealthProbe bool            `json:"HealthProbe"`
	StateChange chan int        `json:"-"`
	ExitAppChan <-chan struct{} `json:"-"`
	GinEngine   *gin.Engine     `json:"-"`
	AppHealthz  bool            `json:"AppHealthz"`
	AppReadyz   bool            `json:"AppReadyz"`
	//private
	exitAppChan    chan struct{} `json:"-"`
	intHealthProbe bool
	intReadyProbe  bool
}

func New(name string) (*Service, error) {
	fmt.Println("New called")
	fmt.Println("New called")

	svc := &Service{
		Name: name,
	}

	return svc, nil
}
