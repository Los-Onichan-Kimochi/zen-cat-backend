FROM --platform=linux/amd64 golang:1.23.7 as base
# Add a work directory
WORKDIR /ASTRO_CAT_BACKEND
# Cache and install dependencies
COPY go.mod go.sum ./
RUN go mod download
# Copy app files
COPY src/ ./src
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /main ./src/server/


FROM gcr.io/distroless/static-debian11
COPY --from=base /main .

# Railway uses PORT environment variable
ENV PORT=8080
EXPOSE 8080

# Health check to ensure the application is ready
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD ["./main", "--health-check"] || exit 1

CMD ["./main"]
