#!/usr/bin/env sh

useradd --system --home /var/lib/pufferpanel --user-group pufferpanel

exitCode=$?
[ $exitCode -ne 0 && $exitCode -ne 9 ] || exit $exitCode