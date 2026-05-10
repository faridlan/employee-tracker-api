# ==========================================
# STAGE 1: BUILDER
# ==========================================
FROM golang:1.26-alpine AS builder

WORKDIR /app

# Copy module files dan download dependencies (caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code
COPY . .

# Build dengan flag tambahan untuk memperkecil ukuran binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o employee-tracker-api ./cmd/api/main.go

# ==========================================
# STAGE 2: RUNNER
# ==========================================
FROM alpine:latest

WORKDIR /app

# Install timezone data dan sertifikat SSL (Penting untuk HTTPS external call)
RUN apk --no-cache add tzdata ca-certificates

# Set default Timezone di level OS Container
ENV TZ=Asia/Jakarta

# Buat non-root user & group demi keamanan (Best Practice)
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Copy binary dari stage builder
COPY --from=builder /app/employee-tracker-api .

# Ubah kepemilikan file binary ke non-root user
RUN chown appuser:appgroup ./employee-tracker-api

# Pindah eksekusi ke user yang baru dibuat
USER appuser

EXPOSE 8080

CMD ["./employee-tracker-api"]