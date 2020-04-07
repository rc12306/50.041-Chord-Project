#!/usr/bin/env bash

# kill all containers
docker rm -f $(sudo docker ps -a -q)

# prune docker systems
docker system prune -af

# clear terminal
clear