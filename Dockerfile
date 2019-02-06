FROM golang:1.11.5-alpine3.9 as builder

RUN mkdir /build 

WORKDIR /build 
COPY go.mod .
COPY go.sum .

RUN apk add git
RUN go mod download
ADD . /build/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o goscout 

FROM scratch
COPY --from=builder /build/goscout /
ENTRYPOINT ["./goscout"]