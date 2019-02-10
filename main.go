package main

import (
    "github.com/undiabler/golang-whois"
    "github.com/imroc/req"
    "github.com/google/logger"
    "encoding/base64"
    "io/ioutil"
    "fmt"
    "strings"
    "os"
)

func main() {
	log := logger.Init("main", true, false, ioutil.Discard)
	log.Info("starting")
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("Need a domain")
		os.Exit(1)
	}
	domain := args[0]
	if strings.Count(domain, ".") > 1 {
		// remove subdomain
		subdomainIndex := strings.Index(domain, ".")
		domain = domain[subdomainIndex+1:]
	}
	result, err := whois.GetWhois(domain)

	if err != nil || len(whois.ParseDomainStatus(result)) == 0 {
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
	}
	log.Info("done.")
}
