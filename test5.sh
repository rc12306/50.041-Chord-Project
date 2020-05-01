#!/usr/bin/env bash

# rebuild
docker-compose build
docker-compose up -d --scale node=5

# this script runs Test1 in 3 Docker containers
# extract all running container id
nodes=$(bash -c "docker ps -a -q")

# bash each container in new window
for node in $nodes
do 
    gnome-terminal --command "bash -c 'docker exec -it $node go test -run Test5; $SHELL'"
done
