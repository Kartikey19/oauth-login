FROM golang:1.12

ENV PORT=8080
WORKDIR /app/server
COPY . . 
RUN go build 
CMD ["./server"]