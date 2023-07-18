FROM golang:1.20

RUN apt-get update \
    && apt-get install -y \
    wkhtmltopdf

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /figs

EXPOSE 8080

CMD ["/figs", "--port", "8080"]