ARG HOST
ARG PORT
ARG SECRET_KEY
ARG DB_USER
ARG DB_NAME
ARG SSL_MODE
ARG DB_PORT
ARG DB_PASS
ARG DB_HOST

FROM golang:alpine AS builder

WORKDIR /build

ADD go.mod .

COPY . .

RUN go build -v -o main

FROM alpine

ENV HOST ${HOST}
ENV PORT ${PORT}

ENV SECRET_KEY ${SECRET_KEY}
ENV DB_USER ${DB_USER}
ENV DB_NAME ${DB_NAME}
ENV SSL_MODE ${SSL_MODE}
ENV DB_PORT ${DB_PORT}
ENV DB_PASS ${DB_PASS}
ENV DB_HOST ${DB_HOST}


ADD .env .
ADD /logs/auth.log /logs/auth.log

COPY --from=builder /build/main /build/main

EXPOSE ${PORT}

CMD ["/build/main"]