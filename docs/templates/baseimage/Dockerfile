# change FROM to output of basegodel (typically godeltutorial:setup) to use warm cache (skips expensive downloads)
FROM golang:1.10.3

ENV GODEL_VERSION 2.9.1
ENV GODEL_CHECKSUM 240052b05e96e95b3f9bae3f89a091ba5dc5ec808d6a8d9cf086be92f7cdd31c

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
