package main

import (
    "github.com/undiabler/golang-whois"
    log "github.com/sirupsen/logrus"
    "github.com/aws/aws-lambda-go/lambda"
    "context"
    "strings"
    "os"
    "strconv"
    "errors"
)
type Event struct {
	Domain string `json:"domain"`
}
type Response struct {
	success bool
	message string
}
var is_lambda bool
var verbose bool
func init() {
	// setup environment variables
	is_lambda = false
	ISLAMBDA := strings.ToLower(os.Getenv("ISLAMBDA"))
	if len(ISLAMBDA) > 0 {
		is_lambda, _ = strconv.ParseBool(ISLAMBDA)
	}

	verbose = false
	VERBOSE := strings.ToLower(os.Getenv("VERBOSE"))
	if len(VERBOSE) > 0 {
		verbose, _ = strconv.ParseBool(VERBOSE)
	}
}
func init() {
	// setup logging
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetOutput(os.Stdout)
	if verbose {
		log.SetLevel(log.InfoLevel)
	} else {
		log.SetLevel(log.ErrorLevel)
	}
}
func IsAvailable(domain string) (bool, error) {
	log.Trace("inside IsAvailable")
	if strings.Count(domain, ".") > 1 {
		log.Info("need to remove the subdomain")
		subdomainIndex := strings.Index(domain, ".")
		domain = domain[subdomainIndex+1:]
	}
	log.Debug("domain = " + domain)
	if len(domain) > 0 {
		result, err := whois.GetWhois(domain)
		if err != nil || len(whois.ParseDomainStatus(result)) == 0 {
			log.Info(domain+" is available")
			return true, nil
		} else {
			log.Warning("got a response. "+domain+" is not available.")
			return false, nil
		}
	} else {
		log.Error("no domain provided")
		return false, errors.New("no domain provided")
	}
}
func LambdaHandler(ctx context.Context, req Event) (Response, error) {
	log.Trace("in LambdaHandler")
	is_available, err := IsAvailable(req.Domain)
	if err != nil {
		return Response{success: false, message: err.Error()}, err
	} else if is_available {
		return Response{success: true, message: req.Domain+" is available"}, nil
	} else {
		return Response{success: true, message: "got a response. "+req.Domain+" is not available."}, nil
	}
}
func CommandLine() {
	log.Trace("in CommandLine")
	args := os.Args[1:]
	if len(args) < 1 {
		log.Error("Need a domain")
		os.Exit(1)
	}
	domain := args[0]
	is_available, _ := IsAvailable(domain)
	if is_available {
		// domain is available
		os.Exit(0)
	} else {
		// domain is unavailable
		os.Exit(2)
	}
}
func main() {
	log.Info("starting")
	if is_lambda {
		// run in lambda function
		lambda.Start(LambdaHandler)
	} else {
		CommandLine()
	}
}
