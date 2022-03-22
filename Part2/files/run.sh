#!/bin/sh
echo "The secret password is $SECRETNAME"
sleep 1
echo "The date and time is $(date)"
sleep 1
echo "Here is a random number $RANDOM "
sleep 1
echo "This is $VAR1 $VAR2!"
sleep 1
echo "Lets list the files"
ls /
if [ -f /keepout/secrets ]
then
  ls /keepout
  echo "The contents of the secret file at /secrets are"
  cat /keepout/secrets
fi