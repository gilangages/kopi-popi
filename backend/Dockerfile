FROM golang:1.26.3-alpine AS builder

# Set working directory di dalam container
WORKDIR /app

# Salin file modul terlebih dahulu (untuk caching layer docker)
COPY go.mod go.sum ./
RUN go mod download

# Salin semua kode sumber
COPY . .

# Build aplikasi Go menjadi file binary statis
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/kopi-popi/main.go

# Menggunakan image scratch atau alpine kosong untuk runtime agar ukurannya kecil
FROM alpine:latest

WORKDIR /app

# Salin binary file yang sudah dibuild dari tahap 'builder'
COPY --from=builder /server .

# Expose port yang digunakan
EXPOSE 8080

# Jalankan server
CMD ["./server"]
