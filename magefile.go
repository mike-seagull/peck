// +build mage

package main

import (
	"fmt"
	"github.com/magefile/mage/sh"
	"github.com/magefile/mage/mg" // mg contains helpful utility functions, like Deps
)
//runs "GOOS=linux go build main.go
func buildLinux() error {
	fmt.Println("building for linux")
	env := map[string]string{
		"GOOS": "linux",
	}
	return sh.RunWith(env, "go", "build", "main.go")
}
// zip up the binary to function.zip
func zip() error {
	fmt.Println("zipping binary to function.zip")
	return sh.Run("zip", "function.zip", "main")
}
// install dependencies
func InstallDeps()  {
	sh.Run("go", "get", "github.com/undiabler/golang-whois")
	sh.Run("go", "get", "github.com/sirupsen/logrus")
	sh.Run("go", "get", "github.com/aws/aws-lambda-go/lambda")
}
// package for lambda
func Package() {
	mg.SerialDeps(buildLinux, zip)
}

// Clean up after yourself
func Clean() {
	build_artifacts := [2]string{"main", "function.zip"}
	for _, artifact := range build_artifacts {
		fmt.Println("removing " + artifact)
		sh.Rm(artifact)
	}
}
