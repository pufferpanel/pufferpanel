#!/usr/bin/env sh

useradd --system --home /var/lib/pufferpanel --user-group pufferpanel

exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode
