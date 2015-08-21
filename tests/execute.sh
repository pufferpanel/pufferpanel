npm run-script test-eslint
eslintRes=$?
npm run-script test-mocha
mochaRes=$?
if [ $mochaRes -ne 0 ] || [ $eslintRes -ne 0 ] ;
then
    exit 1
fi