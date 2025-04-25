FROM golang:1.24.0-bullseye as builder

LABEL org.opencontainers.image.authors="hongyu"

ARG SOURCE_FILES
WORKDIR /app

COPY . /app
RUN cp $SOURCE_FILES ./ && \
  make build SOURCE_FILES=./main.go

FROM busybox:1.34.0-glibc

COPY --from=builder /app/build/ /app

WORKDIR /app

RUN chmod +x main

CMD ["./main"]