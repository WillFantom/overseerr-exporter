ARG GO_VERSION=1.16
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /src
ARG EXPORTER_VERSION=
RUN test -n "$EXPORTER_VERSION"
COPY ./go.mod ./go.mod
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/willfantom/overseerr-exporter/cmd.version=${EXPORTER_VERSION}" -o overseerr-exporter .

FROM alpine

RUN apk --no-cache add ca-certificates
LABEL maintainer="Will Fantom <willf@ntom.dev>"

WORKDIR /exporter
COPY --from=builder /src/overseerr-exporter .
ENTRYPOINT [ "./overseerr-exporter" ]
