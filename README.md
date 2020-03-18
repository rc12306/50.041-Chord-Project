# Chord Implementation on Docker

**Command Line Interface**
![](./CLI.PNG)

**Installing Docker Compose (Linux):**
1. Run this command to download Docker Compose:
```
sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
```
2. Apply executable permissions:
```
sudo chmod +x /usr/local/bin/docker-compose
```
*Installation instructions from [Docker Compose](https://docs.docker.com/compose/install/)*

**Setting up nodes:**

I don't know what these are for but you need these:
```
docker-compose pull
```
```
docker-compose build
```

Run the compose file:
```
docker-compose up -d --scale node=5
```
This creates 1 root node, and 5 other nodes.

To view all containers, run this command:
```
docker ps
```

To view the IP addresses of containers, run this command:
```
docker inspect -f '{{.Name}} - {{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' $(docker ps -aq)
```
