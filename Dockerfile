# Load golang image to build
FROM golang:1.18-alpine as builder
ARG APP_PATH

RUN mkdir -p /app
WORKDIR /app
COPY . .

RUN go mod download
RUN go build -buildvcs=false -o=appbin $APP_PATH


# Deploy execute file to simple linux server
FROM alpine
RUN mkdir -p /app
WORKDIR /app
COPY --chown=0:0 --from=builder /app/ ./

EXPOSE 8080
# Resevre port
EXPOSE 8090

ENTRYPOINT ["/app/appbin"]