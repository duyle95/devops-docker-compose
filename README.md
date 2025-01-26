# How to start

- ci without integration test

docker-compose -p ci -f docker-compose.yml build

docker-compose -p ci -f docker-compose.yml up

- ci with integration test

docker-compose -p ci -f docker-compose.yml -f docker-compose.tests.yml build

docker-compose -p ci -f docker-compose.yml up -d
docker-compose -p ci -f docker-compose.tests.yml up

## Run gitlab runner
gitlab-runner run-single -c ~/.gitlab-runner/config.toml -r devops-final-project

## Testing locally with curl
curl -u admin:admin localhost:8197/state -X PUT -d "RUNNING" \
    -H "Content-Type: text/plain" \
    -H "Accept: text/plain"

curl localhost:8197/state -X PUT -d "PAUSED" \
    -H "Content-Type: text/plain" \
    -H "Accept: text/plain"

curl localhost:8197/state -X PUT -d "INIT" \
    -H "Content-Type: text/plain" \
    -H "Accept: text/plain"

curl localhost:8197/state -X PUT -d "SHUTDOWN" \
    -H "Content-Type: text/plain" \
    -H "Accept: text/plain"

curl localhost:8197/run-log -H "Content-Type: text/plain" -H "Accept: text/plain"

curl localhost:8197/state -X GET -H "Content-Type: text/plain" -H "Accept: text/plain"

curl localhost:8197/request -X GET -H "Content-Type: text/plain" -H "Accept: text/plain"