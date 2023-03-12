# ----------------------------------------------------------------------------------------
# Image: Builder
# ----------------------------------------------------------------------------------------
FROM golang:alpine as builder

# setup the environment
ENV TZ=Europe/Berlin

# install dependencies
RUN apk --update --no-cache add git gcc musl-dev
WORKDIR /work
ADD ./ ./

# build the go binary
RUN go build -ldflags \
        '-X "main.BuildTime='$(date -Iminutes)'" \
         -X "main.GitCommit='$(git rev-parse --short HEAD)'" \
         -X "main.GitBranch='$(git rev-parse --abbrev-ref HEAD)'" \
         -X "main.BuildNumber='$CI_BUILDNR'" \
         -s -w' \
         -v -o /tmp/nsq-discovery-consul .

# ----------------------------------------------------------------------------------------
# Image: Deployment
# ----------------------------------------------------------------------------------------
FROM alpine:latest
MAINTAINER Maximilian Pachl <m@ximilian.info>

# setup the environment
ENV TZ=Europe/Berlin

RUN apk --update --no-cache add ca-certificates tzdata bash su-exec

# add relevant files to container
COPY --from=builder /tmp/nsq-discovery-consul /usr/sbin/nsq-discovery-consul
ADD entry.sh /entry.sh

ENTRYPOINT ["/entry.sh"]
CMD /usr/sbin/nsq-discovery-consul
