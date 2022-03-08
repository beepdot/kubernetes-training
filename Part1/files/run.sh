#!/bin/sh
while true
do
  if [ -f /envconfig ]
  then
    export $(cat /envconfig)
  fi
  echo "Hello $FIRSTNAME $LASTNAME"
  sleep 2
done