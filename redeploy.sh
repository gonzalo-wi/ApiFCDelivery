#!/bin/bash
# Script para redesplegar en producción

echo "🔄 Actualizando código desde Git..."
git pull origin main

echo "🏗️ Reconstruyendo imagen Docker..."
docker build -t gofricalor-api:latest .

echo "🛑 Deteniendo contenedor actual..."
docker stop gofrocalor-api

echo "🗑️ Eliminando contenedor anterior..."
docker rm gofrocalor-api

echo "🚀 Iniciando nuevo contenedor..."
docker run -d \
  --name gofrocalor-api \
  --restart unless-stopped \
  -p 8080:8080 \
  --env-file .env \
  gofricalor-api:latest

echo "✅ Redespliegue completado"
echo ""
echo "📋 Verificando logs..."
docker logs --tail 50 gofrocalor-api
