domain_pecker
===
Checks a domain to see if it is available
___
[![Build Status](https://travis-ci.com/mike-seagull/domain_pecker.svg?branch=master)](https://travis-ci.com/mike-seagull/domain_pecker)

![alt text](woodpecker.png "Woodpecker")
#### Environment Variables
* HOME_API_USER
* HOME_API_AUTH
* HOME_API_DOMAIN
#### To check a domain
```domain_pecker $domain```

if the environment variables are provided, it will send a pushover notification if the domain is available.
