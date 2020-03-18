FROM golang:1.13

# Set current working directory
# Created a new folder called app
WORKDIR /app

# Add files from the current directory to the working directory
ADD . .

# Commands you can run when creating a docker
# Added nano to edit files on docker
RUN apt-get update
RUN apt-get install nano -y
RUN go get -u golang.org/x/sync
