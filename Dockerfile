FROM golang:1.17-alpine AS build

WORKDIR /src/
COPY main.go /src/
COPY go.* /src/
COPY pkg /src/pkg
RUN  CGO_ENABLED=0 go build -o /bin/server

FROM scratch
COPY --from=build /bin/server /bin/server
EXPOSE 8088
ENTRYPOINT ["/bin/server"]
