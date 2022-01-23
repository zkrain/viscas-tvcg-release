FROM golang:1.17-alpine
WORKDIR /backend/src
COPY . .
RUN go build main/main-http-new.go
EXPOSE 8888
ENTRYPOINT ["/backend/src/main-http-new"]

