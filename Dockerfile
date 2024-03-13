# Etapa de construcción
FROM golang:1.21-alpine AS build

WORKDIR /app

# Copia los archivos de descripción de dependencias
COPY go.mod go.sum ./
RUN go mod download

# Copia todo el código fuente
COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -o app .

# Etapa de producción
FROM gcr.io/distroless/static:nonroot

WORKDIR /app

# Copia el binario de la etapa de construcción
COPY --from=build /app/app .

# Expone el puerto
EXPOSE 8080

# Define el comando de entrada
ENTRYPOINT ["./app"]