FROM golang:1.12.4 as build

ENV GO111MODULE on
ENV CGO_ENABLED 0 

COPY . /go/src/github.com/MicroServicesAIOps/HelloWorld-gokit-mongodb
WORKDIR /go/src/github.com/MicroServicesAIOps/HelloWorld-gokit-mongodb

RUN go mod tidy

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app main.go

FROM alpine:3.4

ENV	SERVICE_USER=myuser \
	SERVICE_UID=10001 \
	SERVICE_GROUP=mygroup \
	SERVICE_GID=10001

ENV MONGO_HOST mytestdb:27017
ENV HATEAOS user
ENV USER_DATABASE mongodb

RUN	addgroup -g ${SERVICE_GID} ${SERVICE_GROUP} && \
	adduser -g "${SERVICE_NAME} user" -D -H -G ${SERVICE_GROUP} -s /sbin/nologin -u ${SERVICE_UID} ${SERVICE_USER}

WORKDIR /
COPY --from=0 /app /app

RUN	chmod +x /app && \
	chown -R ${SERVICE_USER}:${SERVICE_GROUP} /app 

USER ${SERVICE_USER}

CMD ["/app"]
EXPOSE 8084