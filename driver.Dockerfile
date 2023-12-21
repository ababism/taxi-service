FROM golang:1.21.4 as build
WORKDIR /app

#COPY go.mod go.sum ./
#RUN go mod vendor
#RUN go mod download

#COPY projects/driver .
COPY . .

#RUN go version
#ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o driver-svc ./projects/driver/cmd/server

## install psql
#RUN apt-get update
#RUN apt-get -y install postgresql-client

FROM alpine:latest as production
COPY --from=build /app/driver-svc ./
#COPY --from=build /app/.env ./
COPY --from=build /app/projects/driver/migrations ./migrations
COPY --from=build /app/projects/driver/config/config.docker.yml ./config/config.local.yml

CMD ["./driver-svc"]
#ENTRYPOINT ["driver-svc"]

EXPOSE 8080