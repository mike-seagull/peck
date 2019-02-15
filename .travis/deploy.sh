#!/usr/bin/env bash
scp -o LogLevel=quiet ./bin/domain_pecker ${REMOTE_USER}@${REMOTE_SERVER}:/usr/local/bin/domain_pecker
ssh -o LogLevel=quiet ${REMOTE_USER}@${REMOTE_SERVER} "/bin/chmod +x /usr/local/bin/domain_pecker"
