FROM golang:1.6
RUN mkdir /files \
    && mkdir -p $GOPATH/src/github.com/feedhenry/openshift-template-tool
WORKDIR $GOPATH/src/github.com/feedhenry/openshift-template-tool
ADD . .
RUN go install
VOLUME /files
ENTRYPOINT ["openshift-template-tool"]