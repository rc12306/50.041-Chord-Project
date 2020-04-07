#!/bin/bash

# open Chord containers
for i in 1 2 3
do
  gnome-terminal --command "bash -c 'docker exec -it 50041chordproject_node_$i bash; $SHELL'"
done