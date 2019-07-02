FROM golang:1.12.5

ARG MAKE_TARGET

ENV EXECUTABLE=wasgeit-${MAKE_TARGET}

ENV GO111MODULE=on

WORKDIR /go/src/github.com/bjorm/wasgeit

ADD . .

RUN make ${MAKE_TARGET}

WORKDIR /wasgeit

USER nobody

ENTRYPOINT ["/go/src/github.com/bjorm/wasgeit/entrypoint.sh"]
