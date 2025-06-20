# Multi-stage build para optimizar el tamaño final de la imagen
FROM golang:1.23-alpine AS builder

# Instalar herramientas necesarias
RUN apk add --no-cache git ca-certificates tzdata

# Crear directorio de trabajo
WORKDIR /app

# Copiar archivos de dependencias primero (mejor uso de cache)
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Build de la aplicación con optimizaciones
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build \
    -ldflags='-w -s -extldflags "-static"' \
    -a -installsuffix cgo \
    -o astrocat-backend \
    ./src/server/main.go

# Segunda etapa: imagen final minimalista
FROM alpine:latest

# Instalar ca-certificates para HTTPS, timezone data y wget para health check
RUN apk --no-cache add ca-certificates tzdata wget

# Crear usuario no-root para seguridad
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Crear directorios necesarios
WORKDIR /app
RUN mkdir -p /app/logs && \
    chown -R appuser:appgroup /app

# Copiar el binario desde la etapa de build
COPY --from=builder /app/astrocat-backend /app/astrocat-backend

# Dar permisos de ejecución y cambiar ownership
RUN chmod +x /app/astrocat-backend && \
    chown -R appuser:appgroup /app

# Cambiar a usuario no-root
USER appuser

# Exponer el puerto que usa la aplicación
EXPOSE 8098

# Variables de entorno por defecto
ENV MAIN_PORT=8098
ENV ENABLE_SQL_LOGS=false
ENV ENABLE_SWAGGER=true

# Health check - usar el endpoint correcto que existe en el backend
HEALTHCHECK --interval=30s --timeout=10s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8098/health-check/ || exit 1

# Comando por defecto
CMD ["/app/astrocat-backend"] 