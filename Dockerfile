# FROM playcourt/jenkins:go1.17-s
FROM golang:1.17.13-alpine

#Set Working Directory
WORKDIR /usr/src/app

COPY ./codebase-http-over-grpc .

#COPY .env ./

# Build Go
RUN go build .

# Expose Application Port
EXPOSE 7000

# Run The Application
CMD ["./go-http-over-grpc"]
