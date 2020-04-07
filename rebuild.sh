#!/bin/bash
docker-compose build
docker-compose up -d --scale node=5