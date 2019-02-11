#!/usr/bin/env bash
scp -o LogLevel=quiet ./bin/whois_checker ${REMOTE_USER}@${REMOTE_SERVER}:/usr/local/bin/whois_checker
ssh -o LogLevel=quiet ${REMOTE_USER}@${REMOTE_SERVER} "/bin/chmod +x /usr/local/bin/whois_checker"
