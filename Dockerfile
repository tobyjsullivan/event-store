FROM golang
ADD . /go/src/github.com/tobyjsullivan/event-store.v3
RUN  go install github.com/tobyjsullivan/event-store.v3

EXPOSE 3000

CMD /go/bin/event-store.v3
