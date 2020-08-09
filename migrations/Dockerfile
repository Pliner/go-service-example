FROM golang as builder
WORKDIR /migrations
RUN go get -u github.com/jackc/tern
COPY . .
CMD ["tern", "migrate"]