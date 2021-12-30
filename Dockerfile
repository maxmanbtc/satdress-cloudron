FROM golang:1.16.0-alpine AS builder

WORKDIR /opt/build

COPY ./*.go ./*.html ./go.mod ./go.sum ./
COPY static ./static

RUN apk add gcc musl-dev linux-headers
RUN go get
RUN go build

FROM alpine:3.14

ENV PORT=17422
ENV DOMAIN=satoshis.se
ENV SECRET=soryakfqomdu3r9d5b1x
ENV SITE_OWNER_URL=https://maxmanbtc.com
ENV SITE_OWNER_NAME=@maxmanbtc
ENV SITE_NAME=Satdress

COPY --from=build /opt/build/satdress /usr/local/bin/

EXPOSE 17422

CMD ["satdress"]