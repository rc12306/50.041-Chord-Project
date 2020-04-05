#!/usr/bin/env bash

# this script runs Test1 in 5 Docker container
for i in 1 2 3 4 5 
do
  gnome-terminal --command "bash -c 'docker exec -it 50041-chord-project_node_$i go test -run Test1; $SHELL'"
done
