# How to start

docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml build

docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml up -d
