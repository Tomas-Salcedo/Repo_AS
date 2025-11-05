# Etapa 1: Build
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Instalar dependencias necesarias
RUN apk add --no-cache git

# Copiar archivos de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o pedido-service .

# Etapa 2: Runtime
FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/pedido-service .

# Exponer puerto
EXPOSE 8083

# Ejecutar
CMD ["./pedido-service"]