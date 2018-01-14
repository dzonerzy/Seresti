#!/usr/bin/env bash

# Seresti processes service
# Each defined variable will be received as an enviroment variable , the prefix
# is SERESTI_ , so if variable is named person then GRESTY_PERSON is out variable

if [ -z "${SERESTI_FILTER}" ]; then
    ps uax | awk -v SEP="$SERESTISEP" '{print $1 SEP $2 SEP $11}' | tail -n +2
else
    ps uax | awk -v SEP="$SERESTISEP"  '{print $1 SEP $2 SEP $11}' | grep $SERESTI_FILTER
fi
