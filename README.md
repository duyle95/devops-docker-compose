# How to start

- ci without integration test

docker-compose -p ci -f docker-compose.yml build

docker-compose -p ci -f docker-compose.yml up -d

- ci with integration test

docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml build

docker-compose -p ci -f docker-compose.yml up -d
docker-compose -p ci -f docker-compose.tests.yml up
