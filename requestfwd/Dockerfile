FROM golang:1.11 as builder

RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -u github.com/gorilla/mux && \
    go get -u github.com/jinzhu/copier && \
    go get -u github.com/ttacon/glog

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o requestfwd requestfwd.go
FROM scratch
COPY --from=builder /build/requestfwd /
WORKDIR /
CMD ["./requestfwd"]