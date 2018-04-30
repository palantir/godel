# change FROM to output of basegodel to use warm cache (skips expensive downloads)
FROM godeltutorial:setup

ENV GODEL_VERSION 2.0.0-rc12
ENV GODEL_CHECKSUM 6645d041f2243b146f471f79f80fa59adabc5a4a9f0ae9eb1023c3ec9dcac116

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
