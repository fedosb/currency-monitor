FROM golang:1.22.3-alpine AS build
WORKDIR /src

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

RUN go build -a -installsuffix cgo -o bin/app ./cmd/main.go

FROM alpine AS bin

COPY --from=build /src/bin/app /

EXPOSE 8081
CMD ["./app"]