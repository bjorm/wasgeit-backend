FROM golang:1.12.5

ARG MAKE_TARGET

ENV GO111MODULE=on

WORKDIR /go/src/github.com/bjorm/wasgeit

ADD . .

RUN make ${MAKE_TARGET}

WORKDIR /wasgeit

RUN echo wasgeit-${MAKE_TARGET} > run.sh && chmod +x run.sh

USER nobody
CMD ["sh", "-c", "./run.sh"]