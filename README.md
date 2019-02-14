[![Build Status](https://travis-ci.com/mike-seagull/whois_checker.svg?branch=master)](https://travis-ci.com/mike-seagull/whois_checker)
<p>Checks a domain to see if it is available<p/>
<h4>Environment Variables</h4>

* HOME_API_USER
* HOME_API_AUTH
* HOME_API_DOMAIN

<h5>To check a domain</h5>
<code>domain_pecker $domain</code>
<p>if the environment variables are provided, it will send a pushover notification if the domain is available.</p>
