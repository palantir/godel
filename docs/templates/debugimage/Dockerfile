# change FROM to output of basegodel to use warm cache (skips expensive downloads)
FROM godeltutorial:Tutorial

ARG GODEL_ARTIFACT_NAME

ADD ${GODEL_ARTIFACT_NAME} /downloads/
ADD repository /m2/repository

RUN cp -r /downloads/"$(ls /downloads/)"/wrapper/* . && \
    mkdir -p ~/.godel/dists/ && \
    mv /downloads/"$(ls /downloads/)" ~/.godel/dists/ && \
    rm -rf /downloads && \
    ./godelw version && \
    git add godel godelw && \
    git commit -m "Add godel to project"
