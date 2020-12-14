# Dockerfile.deploy

FROM golang:1.15 as builder

ENV APP_USER app
ENV APP_HOME /go/src/cmm

RUN groupadd $APP_USER && useradd -m -g $APP_USER -l $APP_USER
RUN mkdir -p $APP_HOME && chown -R $APP_USER:$APP_USER $APP_HOME

WORKDIR $APP_HOME
USER $APP_USER
COPY go.* ./
COPY main.go .
COPY client client

ENV GO111MODULE on

RUN go mod download
RUN go mod verify
RUN go build -o cmm

FROM golang:1.15-alpine

ENV APP_USER app
ENV APP_HOME /go/src/cmm

RUN mkdir -p $APP_HOME
WORKDIR $APP_HOME

COPY --chown=0:0 --from=builder $APP_HOME/cmm $APP_HOME

EXPOSE 8888
CMD ["./cmm"]