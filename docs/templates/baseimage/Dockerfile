# change FROM to output of basegodel to use warm cache (skips expensive downloads)
FROM godeltutorial:setup

ENV GODEL_VERSION 2.4.0
ENV GODEL_CHECKSUM 48e50289db0acfce7f8a68b12f8ab433a2fcbf67e80b7658dfc1e105228228c1

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
