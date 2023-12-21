FROM golang:1.21.4 as build
WORKDIR /app

#COPY go.mod go.sum ./
#RUN go mod vendor
#RUN go mod download

#COPY projects/location .
COPY . .

#RUN go version
#ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o location-svc ./projects/location/cmd/server

## install psql
#RUN apt-get update
#RUN apt-get -y install postgresql-client

FROM alpine:latest as production
COPY --from=build /app/location-svc ./
#COPY --from=build /app/.env ./
COPY --from=build /app/projects/location/config/config.docker.yml ./config/config.local.yml

CMD ["./location-svc"]
#ENTRYPOINT ["location-svc"]

EXPOSE 8080