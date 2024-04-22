package servicebase

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"
)

type Option func(*Service)

func WithPort(port int) Option {
	return func(svc *Service) {
		svc.Port = port
	}
}

func WithHealthProbe(healthProbe bool) Option {
	return func(svc *Service) {
		svc.HealthProbe = healthProbe
	}
}

func WithVersion(version int) Option {
	return func(svc *Service) {
		svc.Version = version
	}
}

type Service struct {
	//public
	Name        string          `json:"Name"`
	Port        int             `json:"Port"`
	HealthProbe bool            `json:"HealthProbe"`
	ExitAppChan <-chan struct{} `json:"-"`
	GinEngine   *gin.Engine     `json:"-"`
	AppHealthz  bool            `json:"AppHealthz"`
	AppReadyz   bool            `json:"AppReadyz"`
	Version     int             `json:"Version"`
	//private
	exitAppChan    chan struct{} `json:"-"`
	intHealthProbe bool
	intReadyProbe  bool
}

func New(name string, options ...Option) (*Service, error) {
	fmt.Println("+ New Service... " + name)
	exitChan := make(chan struct{})

	r := gin.Default()

	svc := &Service{
		Name:           name,
		Port:           8080,
		HealthProbe:    false,
		ExitAppChan:    exitChan,
		GinEngine:      r,
		AppHealthz:     false,
		AppReadyz:      false,
		exitAppChan:    exitChan,
		intHealthProbe: false,
		intReadyProbe:  false,
	}

	for _, option := range options {
		option(svc)
	}

	//TODO: For now we default to version 1
	if svc.Version == 0 {
		svc.Version = 1
	}

	return svc, nil
}

func Start(svc *Service) {
	defer fmt.Println("Exiting Start...")
	exitApp := make(chan struct{})

	var wg sync.WaitGroup

	srv := &http.Server{
		Addr:    ":8080",
		Handler: svc.GinEngine,
	}

	setupInternalRoutes(svc)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listening: %s\n", err)
		}
	}()

	if svc.HealthProbe {
		wg.Add(1)
		go kubeProbes(9999, svc, exitApp, &wg)
	}

	sig := make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIPE)

	//declare we are ready
	svc.intReadyProbe = true
	svc.intHealthProbe = true

	<-sig
	fmt.Println("Received one of these: -> os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGPIP")

	//TODO: i need to merge the context cancel context stuff this looks ugly and a bit brittle
	close(exitApp)
	wg.Wait()

	closeServer(srv, 3*time.Second)

	close(svc.exitAppChan)
}

func closeServer(srv *http.Server, timeout time.Duration) {

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server shutdown: ", err)
	}

	select {
	case <-ctx.Done():
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Shutdown timed out")
		}
	}

	log.Println("Server exiting")

}

func setupInternalRoutes(svc *Service) {
	svc.GinEngine.GET("/servicebase", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "hiya"})
	})

	svc.GinEngine.GET("/servicebase/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, &svc)
	})
}

func kubeProbes(workerId int, svc *Service, stopChan <-chan struct{}, wg *sync.WaitGroup) {
	workerid := strconv.Itoa(workerId)

	defer fmt.Println("worker finished: " + workerid)
	defer wg.Done()

	fmt.Println("Starting worker: " + workerid)

	svc.GinEngine.GET("/healthz", func(c *gin.Context) {
		if svc.intHealthProbe && svc.AppHealthz {
			c.JSON(http.StatusOK, gin.H{"STATUS": "OK"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"STATUS": "ERROR"})
		}
	})

	svc.GinEngine.GET("/readyz", func(c *gin.Context) {
		if svc.intReadyProbe && svc.AppReadyz {
			c.JSON(http.StatusOK, gin.H{"STATUS": "OK"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"STATUS": "ERROR"})
		}
	})

	<-stopChan
}
