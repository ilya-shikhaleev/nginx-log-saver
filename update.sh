docker build -t nginx-log-saver .
docker stop ilya-container
docker rm ilya-container
docker run --publish 6060:8080 --name ilya-container --link mongo-container --rm nginx-log-saver