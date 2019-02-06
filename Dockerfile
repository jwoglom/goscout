FROM golang:1.11 as builder

RUN mkdir /build 
ADD . /build/
WORKDIR /build 
RUN go get -u github.com/go-gorp/gorp && \
    go get -u github.com/gorilla/mux && \
    go get -u github.com/jinzhu/copier && \
    go get -u github.com/mattn/go-sqlite3 && \
    go get -u github.com/ttacon/glog

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o goscout goscout.go
FROM scratch
COPY --from=builder /build/goscout /
WORKDIR /
CMD ["./goscout"]