#Initial Stage

FROM golang:1.14.2-alpine3.11 as build-env

RUN apk add --no-cache git
RUN apk update && apk upgrade libcurl && apk add git openssh-client curl gcc
#RUN git config --global http.https://gopkg.in.followRedirects true
#RUN curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $(go env GOPATH)/bin v1.17.1
# Set the Current Working Directory inside the container, to enable module features
ENV GO111MODULE on

ENV WKDIR /app
WORKDIR ${WKDIR}
# Copy go mod and sum files
COPY go.mod ${WKDIR}
COPY go.sum ${WKDIR}
COPY *.go ${WKDIR}/

COPY config ${WKDIR}/config
COPY productService ${WKDIR}/productService
COPY tests ${WKDIR}/tests
COPY xeroErrors ${WKDIR}/xeroErrors
COPY xeroHelper ${WKDIR}/xeroHelper
COPY xeroLog ${WKDIR}/xeroLog
RUN mkdir -p  ${WKDIR}/cert/
COPY ./_output/cert ${WKDIR}/cert
RUN mkdir -p  ${WKDIR}/data/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

#RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 golangci-lint run ./... --issues-exit-code=0 --no-config --deadline=3m --disable-all --enable=deadcode --enable=structcheck --enable=typecheck --enable=unused --enable=varcheck --enable=goimports
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -i -o ${WKDIR}/xeroProductAPI



#Final Stage
FROM alpine:3.1

#create and run go program as seperate user
RUN adduser -D -u 10000 gouser
USER gouser
WORKDIR /
COPY --from=build-env /app /
#ENV SERVICE_ADDR :8888
EXPOSE 8080
EXPOSE 8081

CMD ["/xeroProductAPI"]