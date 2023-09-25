# Build the testserver binary
# Docker command must be invoked from the projec root directory
FROM golang:1.20 as builder

ARG TARGETARCH

ENV GOARCH=$TARGETARCH

WORKDIR /src

# Copy the go manifests and source
COPY . .

# Build
RUN go build -o pingserver ./main.go

FROM ubuntu:latest

WORKDIR /
COPY --from=builder /src/pingserver .
USER 0:0

CMD [ "/pingserver" ]