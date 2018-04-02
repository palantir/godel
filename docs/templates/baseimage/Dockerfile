# change FROM to output of basegodel to use warm cache (skips expensive downloads)
FROM godeltutorial:setup

ENV GODEL_VERSION 2.0.0-rc8
ENV GODEL_CHECKSUM 534e2a5c22ac113cf025cacbe748a6f031deb26e8c97cd14d89df3ef20f57595

ENV PROJECT_PATH github.com/nmiyake/echgo2

ENV GIT_USERNAME "Tutorial User"
ENV GIT_EMAIL "tutorial@tutorial-user.com"

RUN apt-get update && apt-get install -y tree

# Set up Git author parameters and create initial repository directory
RUN git config --global user.name "${GIT_USERNAME}" && \
    git config --global user.email "${GIT_EMAIL}" && \
    mkdir -p ${GOPATH}/src/${PROJECT_PATH}

WORKDIR /go/src/${PROJECT_PATH}
