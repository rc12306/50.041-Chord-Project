# Chord Implementation on Docker
**Installing Docker Compose (Linux):**
1. Run this command to download Docker Compose: `sudo curl -L "https://github.com/docker/compose/releases/download/1.25.4/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose`
2. Apply executable permissions: `sudo chmod +x /usr/local/bin/docker-compose`
__Set-up instructions from https://docs.docker.com/compose/install/ __

**Setting up nodes:**
Run these commands:
`docker build .`
`docker-compose up -d`
This builds the image, then sets up a root node to communicate with the client and 1 other node. 
