FROM golang
ADD . /go/src/github.com/tobyjsullivan/event-store.v3
RUN  go install github.com/tobyjsullivan/event-store.v3
CMD /go/bin/event-store.v3
