# Etapa de compilación
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o authApp .

# Etapa final: imagen liviana
FROM alpine:3.18

WORKDIR /app

RUN apk update && apk upgrade --no-cache \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

# Copiar solo el binario final
COPY --from=builder /app/authApp .

# Asignar propietario y permisos (lectura/ejecución sin escritura)
RUN chown appuser:appgroup authApp \
    && chmod 700 authApp

# Cambiar al usuario no root
USER appuser

# Exponer puerto
EXPOSE 8081

# Ejecutar la app
CMD ["./authApp"]