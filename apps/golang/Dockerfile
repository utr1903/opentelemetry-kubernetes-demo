# syntax=docker/dockerfile:1

### Build
FROM golang:1.20-buster AS build

ARG APP_NAME

WORKDIR /app

COPY ./commons ./commons
COPY ./${APP_NAME} ./${APP_NAME}

WORKDIR /app/${APP_NAME}

RUN go mod download

RUN go build -o ./out .

### Deploy
FROM gcr.io/distroless/base-debian10

ARG APP_NAME

WORKDIR /

COPY --from=build /app/${APP_NAME}/out /out

EXPOSE 8080

USER nonroot:nonroot

CMD ["./out"]