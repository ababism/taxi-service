FROM golang:1.21.4 as build
WORKDIR /app

COPY . .

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o location-svc ./projects/location/cmd

FROM alpine:latest as production

COPY --from=build /app/location-svc ./

COPY --from=build /app/.env ./
COPY --from=build /app/projects/location/migrations ./migrations
COPY --from=build /app/projects/location/config/config.docker.yml ./config/config.local.yml

CMD ["./location-svc"]

EXPOSE 8080