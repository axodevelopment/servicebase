package utils

import (
	"fmt"
	"os"
)

type EnvVar struct {
	Value  string `json:"value"`
	Exists bool   `json:"exists"`
}

// Very basic parsing
func GetEnvVars(envvars ...string) map[string]EnvVar {
	vars := make(map[string]EnvVar)

	for _, v := range envvars {
		ev := os.Getenv(v)

		env := EnvVar{ev, ev != ""}

		vars[v] = env

		if ev == "" {
			fmt.Println("OsEnvVar NotFound - [" + v + "] => defaulted to ''")
		}
	}

	return vars
}
