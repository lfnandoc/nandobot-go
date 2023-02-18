FROM alpine

RUN apk --update upgrade
RUN apk add sqlite
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

RUN apk add --no-cache go

RUN rm -rf /var/cache/apk/*

WORKDIR /app

COPY . .

RUN go build -o nandobot-go . 

EXPOSE 8080

CMD ["./nandobot-go"]