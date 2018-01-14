#!/usr/bin/env bash

# Seresti hello service
# Each defined variable will be received as an enviroment variable , the prefix
# is SERESTI_ , so if variable is named person then GRESTY_PERSON is out variable

echo "Hello $SERESTI_PERSON!"
