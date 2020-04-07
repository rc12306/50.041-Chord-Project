#!/bin/bash
docker-compose build
docker-compose up -d --scale node=5

for i in 1 2 3 4 5
do
  gnome-terminal --command "bash -c 'docker exec -it 50041chordproject_node_$i bash; $SHELL'"
done