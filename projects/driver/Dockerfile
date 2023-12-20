FROM golang:1.21-bullseye as build
WORKDIR /bin
COPY . .

#RUN go version
ENV GOPATH=/
ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o user-svc ./cmd/server/main.go

## install psql
#RUN apt-get update
#RUN apt-get -y install postgresql-client

FROM alpine:latest as production
COPY --from=build /bin/user-svc ./
COPY --from=build  /bin/config/dockerconfig.yml ./config/config.yml

CMD ["./user-svc"]
#ENTRYPOINT ["user-svc"]

EXPOSE 8080