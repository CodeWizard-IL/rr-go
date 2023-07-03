FROM golang:1.20 AS builder
ARG PACKAGE

COPY . /app/

RUN find /app

WORKDIR /app/

RUN go mod download

WORKDIR /app/${PACKAGE}

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o ${PACKAGE} .

FROM gcr.io/distroless/static
ARG PACKAGE

USER nonroot:nonroot

COPY --from=builder --chown=nonroot:nonroot /app/${PACKAGE}/${PACKAGE} /app

ENTRYPOINT ["/app"]