#!/bin/sh
docker inspect -f '{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -aq) | tee ip.txt
# this script could be incorporated into the Dockerfile / docker_compose file
# so it will run at startup 
