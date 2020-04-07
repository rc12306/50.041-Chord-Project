#!/usr/bin/env bash

# rebuild
docker-compose build
docker-compose up -d --scale node=3

# this script runs Test1 in 3 Docker containers
for i in 1 2 3
do
  gnome-terminal --command "bash -c 'docker exec -it 50041-chord-project_node_$i go test -run Test1; $SHELL'"
done
