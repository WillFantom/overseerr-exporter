FROM golang:1.19-alpine3.17 as builder

WORKDIR /src

ARG EXPORTER_VERSION=
RUN test -n "$EXPORTER_VERSION"

COPY ./go.mod ./go.mod
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X github.com/willfantom/overseerr-exporter/cmd.version=${EXPORTER_VERSION}" -o overseerr-exporter .


FROM alpine:3.17

RUN apk --no-cache add ca-certificates
LABEL maintainer="Will Fantom <willf@ntom.dev>"

WORKDIR /exporter
COPY --from=builder /src/overseerr-exporter .

ENTRYPOINT [ "./overseerr-exporter" ]
