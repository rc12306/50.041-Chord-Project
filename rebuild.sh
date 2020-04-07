#!/bin/bash
docker-compose build
docker-compose up -d --scale node=5

# extract all running container id
nodes=$(bash -c "docker ps -a -q")

# bash each container in new window
for node in $nodes
do 
    gnome-terminal --command "bash -c 'docker exec -it $node bash; $SHELL'"
done
