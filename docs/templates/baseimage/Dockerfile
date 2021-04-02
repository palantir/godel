# change FROM to output of basegodel (typically godeltutorial:setup) to use warm cache (skips expensive downloads)
FROM golang:1.16.2

ENV GODEL_VERSION 2.36.0
ENV GODEL_CHECKSUM 91137f4fb9e1b4491d6dd821edf6ed39eb66f21410bf645a062f687049c45492

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
