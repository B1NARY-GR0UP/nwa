FROM golang:1.23 AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod tidy

# copy source and build
# disable cgo can reduce the size of binary file
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build .


# make a bare minimal image
FROM scratch

# source to be scanned should be mounted to /src
WORKDIR /src
COPY --from=build /app/nwa /app/nwa

ENTRYPOINT ["/app/nwa"]

LABEL org.opencontainers.image.source=https://github.com/B1NARY-GR0UP/nwa
LABEL org.opencontainers.image.description="NWA - A More Powerful License Header Management Tool"
LABEL org.opencontainers.image.licenses=Apache
