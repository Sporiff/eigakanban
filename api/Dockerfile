FROM golang:1.24-alpine AS go-build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

RUN goose up
RUN CGO_ENABLED=0 GOOS=linux go build -o /eigakanban

FROM golang:1.24-alpine AS go-dev

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download
RUN go install github.com/air-verse/air@latest
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

EXPOSE 8080

CMD ./scripts/start.sh

FROM nginx:alpine AS go-prod

COPY --from=go-build /eigakanban /usr/share/nginx/html/eigakanban

EXPOSE 80

CMD ["nginx", "-g", "daemon off;"]