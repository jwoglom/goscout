FROM golang:1.11.5 as builder

RUN mkdir /build 

WORKDIR /build 
COPY go.mod .
COPY go.sum .

RUN go mod download
ADD . /build/
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o goscout 

FROM scratch
COPY --from=builder /build/goscout /
ENTRYPOINT ["./goscout"]
