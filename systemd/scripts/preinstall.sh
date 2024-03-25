#!/usr/bin/env sh

if [ -n "${CODESPACES}" ]; then
  echo 'CODESPACE IS NOT A PERMITTED PLATFORM!'
  echo 'DO NOT USE CODESPACES TO RUN THIS PACKAGE'
  echo 'THIS IS A VIOLATION OF GITHUB TOS'
  echo 'DO NOT ASK HOW TO GET AROUND IT.'
  exit 1
fi

useradd --system --home /var/lib/pufferpanel --user-group pufferpanel

exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode
