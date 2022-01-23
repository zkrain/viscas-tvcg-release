FROM golang:1.17-alpine
WORKDIR /backend
COPY . .
RUN cd src && go build main/main-http-new.go
WORKDIR /backend/src
EXPOSE 8888
ENTRYPOINT ["/backend/src/main-http-new"]
