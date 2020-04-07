#!/bin/bash
docker-compose build
docker-compose up -d --scale node=5

open Chord containers
for i in 1 2 3 4 5
do
  gnome-terminal --command "bash -c 'docker exec -it 50041chordproject_node_$i go run chord.go; $SHELL'"
done
