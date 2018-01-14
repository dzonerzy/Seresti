#!/usr/bin/env bash

# Seresti info service
# Each defined variable will be received as an enviroment variable , the prefix
# is SERESTI_ , so if variable is named person then GRESTY_PERSON is out variable

echo "$(hostname)$SERESTISEP$(uname)$SERESTISEP$(uptime)"
