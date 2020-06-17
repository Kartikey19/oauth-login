FROM golang:1.13.12-alpine3.11

ENV PORT=8080
WORKDIR /app/server
COPY . . 
RUN go build 
CMD ["./recro_demo"]