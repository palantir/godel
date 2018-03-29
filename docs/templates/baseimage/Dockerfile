# change FROM to output of basegodel to use warm cache (skips expensive downloads)
FROM godeltutorial:setup

ENV GODEL_VERSION 2.0.0-rc5
ENV GODEL_CHECKSUM 08d9ed3e33e69006a9c58ec65cef0ad9bd17af4c73b5c1d1aa116e813a954314

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
