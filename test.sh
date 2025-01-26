RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m'
cleanup () {
  docker-compose -p ci kill
  docker-compose -p ci rm -f
}
trap 'cleanup ; printf "${RED}Tests Failed For Unexpected Reasons${NC}\n"' HUP INT QUIT PIPE TERM
docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml build && docker-compose -p ci -f docker-compose.yml up -d && docker-compose -p ci -f docker-compose.tests.yml up -d

if [ $? -ne 0 ] ; then
  printf "${RED}Docker Compose Failed${NC}\n"
  exit -1
fi

TEST_EXIT_CODE=`docker wait ci-tests-1`
echo "TEST_EXIT_CODE: '$TEST_EXIT_CODE'"

if [ "$TEST_EXIT_CODE" -ne 0 ] ; then
  printf "${RED}Tests Failed${NC} - Exit Code: $TEST_EXIT_CODE\n"
else
  printf "${GREEN}Tests Passed${NC}\n"
fi
cleanup
exit $TEST_EXIT_CODE