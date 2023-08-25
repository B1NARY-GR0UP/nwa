FROM golang:1.21 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# copy source and build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .


# make a bare minimal image
FROM scratch

# source to be scanned should be mounted to /src
WORKDIR /src
COPY --from=build /app/nwa /app/nwa

ENTRYPOINT ["/app/nwa"]