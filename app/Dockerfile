# --- ETAPA 1: EL CONSTRUCTOR (BUILDER) ---
# Usamos una imagen que tiene todas las herramientas de Go (compilador, librerías)
# Es pesada (aprox 800MB)
FROM golang:1.25-alpine AS builder

# Creamos una carpeta de trabajo adentro del contenedor
WORKDIR /app

# 1. Copiamos SOLO los archivos de definición de dependencias
COPY go.mod ./

# 2. Descargamos las librerías de internet (Caché de Docker)
RUN go mod download  

# 3. Copiamos el resto del código fuente (los .go)
COPY . .

# 4. Compilamos el binario
# -o log-gen: El nombre del archivo final será 'log-gen'
RUN go build -o log-gen main.go


# --- ETAPA 2: EL CORREDOR (RUNNER) ---
# Usamos Alpine Linux. Es una distribución minúscula (aprox 5MB).
# No tiene compilador de Go, ni herramientas extra. Solo lo básico para vivir.
FROM alpine:latest

WORKDIR /app

# ¡LA MAGIA!
# Copiamos SOLO el archivo binario 'log-gen' desde la Etapa 1 (builder).
# Tiramos a la basura todo el código fuente y el compilador pesado.
COPY --from=builder /app/log-gen .

# Comando que se ejecuta al iniciar el contenedor
CMD ["./log-gen"]