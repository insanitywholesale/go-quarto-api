# build stage
FROM golang:1.17 as build

ENV CGO_ENABLED 0
ENV GO111MODULE on

WORKDIR /go/src/adise1941

COPY go.mod go.sum ./
RUN go mod download -x

COPY . .

RUN go vet -v
RUN make installwithvars

# run stage
FROM busybox as run

COPY --from=build /go/bin/adise1941 /

CMD ["/adise1941"]
