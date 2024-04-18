package servicebase

import (
	"fmt"
)

type Service struct {
	Name string `json:"Name"`
}

func New(name string) (*Service, error) {
	fmt.Println("New called")
	fmt.Println("New called")

	svc := &Service{
		Name: name,
	}

	return svc, nil
}
