#!/usr/bin/env sh

/pufferpanel/bin/pufferpanel db migrate
exitCode=$?
[ $exitCode -eq 0 ] || [ $exitCode -eq 9 ] || exit $exitCode

exec /pufferpanel/bin/pufferpanel run
