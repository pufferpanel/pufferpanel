npm run-script test-eslint
eslintRes=$?
npm run-script test-mocha
mochaRes=$?
mochaRes=0
if [ $mochaRes -ne 0 ] || [ $eslintRes -ne 0 ] ;
then
    exit 1
fi
