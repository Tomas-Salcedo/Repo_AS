# Etapa de compilación
FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .

# Compilar el binario
RUN CGO_ENABLED=0 GOOS=linux go build -o frontendApp .

# Etapa final: imagen liviana
FROM alpine:3.18

WORKDIR /app

RUN apk update && apk upgrade --no-cache \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

# Copiar el binario y los templates
COPY --from=builder /app/frontendApp .
COPY --from=builder /app/templates ./templates

# Asignar propietario y permisos
RUN chown -R appuser:appgroup . \
    && chmod 700 frontendApp \
    && chmod -R 500 templates

# Cambiar al usuario no root
USER appuser

# Exponer puerto
EXPOSE 8080

# Ejecutar la app
CMD ["./frontendApp"]