FROM golang:1.10.2

ARG MAKE_TARGET

RUN curl -fsSL -o /usr/local/bin/dep https://github.com/golang/dep/releases/download/v0.4.1/dep-linux-amd64
RUN chmod +x /usr/local/bin/dep

WORKDIR /go/src/github.com/bjorm/wasgeit
ADD . .
RUN rm -rf vendor
RUN dep ensure -vendor-only
RUN make ${MAKE_TARGET}

WORKDIR /wasgeit
USER nobody

CMD ["wasgeit-server"]