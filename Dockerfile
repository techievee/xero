#Initial Stage

FROM golang:1.14.2-alpine3.11 as build-env

RUN apk add --no-cache git
RUN apk update && apk upgrade libcurl && apk add git openssh-client curl gcc musl-dev

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
COPY database ${WKDIR}/database
COPY apiServer ${WKDIR}/apiServer
COPY tests ${WKDIR}/tests
COPY xeroErrors ${WKDIR}/xeroErrors
COPY xeroHelper ${WKDIR}/xeroHelper
COPY xeroLog ${WKDIR}/xeroLog
RUN mkdir -p  ${WKDIR}/cert/
COPY ./_output/cert ${WKDIR}/cert
RUN mkdir -p  ${WKDIR}/data/

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

RUN CGO_ENABLED=1 go build -o ${WKDIR}/xeroProductAPI



#Final Stage
FROM alpine:3.1

#create and run go program as seperate user
RUN mkdir -p /data
WORKDIR /
COPY --from=build-env /app /
COPY --from=build-env /app/cert /cert
COPY --from=build-env /app/config /config
#ENV SERVICE_ADDR :8888
EXPOSE 8080
EXPOSE 8081

CMD ["/xeroProductAPI"]