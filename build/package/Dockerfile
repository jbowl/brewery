#docker create -p 50000:50000 --name breweryctr brewery



FROM golang:alpine as build
WORKDIR /app

COPY ./brewery /src

RUN cd /src/cmd/brewery && CGO_ENABLED=0 go build -ldflags="-s -w" -o /app/brewery

FROM alpine
WORKDIR /app
COPY --from=build /app/brewery .

CMD ./brewery

#ENTRYPOINT [ "app/ptaasaiq" ]