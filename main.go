package main

import (
    "github.com/undiabler/golang-whois"
    "github.com/imroc/req"
    "github.com/google/logger"
    "github.com/aws/aws-lambda-go/lambda"
    "encoding/base64"
    "io/ioutil"
    "context"
    "fmt"
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
func IsAvailable(domain string) (bool, error) {
	log := logger.Init("IsAvailable", true, false, ioutil.Discard)
	log.Info("inside IsAvailable")
	if strings.Count(domain, ".") > 1 {
		// remove subdomain
		subdomainIndex := strings.Index(domain, ".")
		domain = domain[subdomainIndex+1:]
	}
	if len(domain) > 0 {
		result, err := whois.GetWhois(domain)
		if err != nil || len(whois.ParseDomainStatus(result)) == 0 {
			log.Info(domain+" is available")
			return true, nil
		} else {
			log.Info("got a response. "+domain+" is not available.")
			return false, nil
		}
	} else {
		return false, errors.New("no domain provided")
	}
}
func LambdaHandler(ctx context.Context, req Event) (Response, error) {
	log := logger.Init("LambdaHandler", true, false, ioutil.Discard)
	fmt.Println(ctx)
	log.Info("in LambdaHandler")
	fmt.Println(req)
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
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Need a domain")
		os.Exit(1)
	}
	domain := args[0]
	log := logger.Init("CommandLine", true, false, ioutil.Discard)
	log.Info("in CommandLine")
	is_available, _ := IsAvailable(domain)
	if is_available {
		// domain is available
		log.Info(domain+" is available")
		home_api_user, user_is_set := os.LookupEnv("HOME_API_USER")
		home_api_auth, auth_is_set := os.LookupEnv("HOME_API_AUTH")
		home_api_domain, domain_is_set := os.LookupEnv("HOME_API_DOMAIN")
		if user_is_set && auth_is_set && domain_is_set {
			// send pushover message via home-api
			log.Info("going to send a push message.")
			user_auth := base64.StdEncoding.EncodeToString([]byte(home_api_user+":"+home_api_auth))
			header := req.Header {
				"Accept":        "application/json",
				"Authorization": "Basic "+user_auth,
			}
			param := req.Param {
				"message": domain+" is available",
				"title":  "whois notification",
			}
			// only url is required, others are optional.
			_, err := req.Post("https://"+home_api_domain+"/api/pushover", header, param)
			if err != nil {
				log.Error(err)
			}
		}
	} else {
		// domain is unavailable
		log.Info("got something back. this domain is in use. going to do nothing.")
		os.Exit(2)
	}
}
func main() {
	log := logger.Init("main", true, false, ioutil.Discard)
	log.Info("starting")
	ISLAMBDA := strings.ToLower(os.Getenv("ISLAMBDA"))
	is_lambda := false
	if len(ISLAMBDA) > 0 {
		is_lambda, _ = strconv.ParseBool(ISLAMBDA)
	}
	if is_lambda {
		// run in lambda function
		lambda.Start(LambdaHandler)
	} else {
		CommandLine()
	}
}
