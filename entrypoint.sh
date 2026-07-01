#!/usr/bin/env sh

/pufferpanel/bin/pufferpanel db upgrade
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

exec /pufferpanel/bin/pufferpanel run
