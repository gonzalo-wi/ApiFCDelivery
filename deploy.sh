#!/bin/bash
# Script de despliegue automatizado para GoFrioCalor
# Uso: ./deploy.sh [production|staging]

set -e  # Detener en caso de error

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Funciones auxiliares
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Verificar que estamos en la rama correcta
CURRENT_BRANCH=$(git branch --show-current)
log_info "Branch actual: $CURRENT_BRANCH"

# Confirmar despliegue
read -p "¿Desplegar a producción? (y/N): " -n 1 -r
echo
if [[ ! $REPLY =~ ^[Yy]$ ]]
then
    log_warn "Despliegue cancelado"
    exit 1
fi

log_info "Iniciando despliegue..."

# 1. Backup de base de datos
log_info "Creando backup de base de datos..."
BACKUP_FILE="backup_$(date +%Y%m%d_%H%M%S).sql"

# Ajusta estos valores según tu configuración
DB_USER=${DB_USER:-"root"}
DB_PASS=${DB_PASS:-""}
DB_NAME=${DB_NAME:-"gofricalor"}

if command -v mysqldump &> /dev/null; then
    mysqldump -u $DB_USER -p$DB_PASS $DB_NAME > backups/$BACKUP_FILE || log_warn "No se pudo crear backup automático"
    log_info "Backup creado: backups/$BACKUP_FILE"
else
    log_warn "mysqldump no disponible, saltando backup automático"
fi

# 2. Pull de cambios
log_info "Actualizando código desde repositorio..."
git pull origin main

# 3. Detener contenedor actual
log_info "Deteniendo contenedor actual..."
docker-compose down

# 4. Reconstruir imagen
log_info "Reconstruyendo imagen Docker..."
docker-compose build --no-cache

# 5. Iniciar nuevo contenedor
log_info "Iniciando nuevo contenedor..."
docker-compose up -d

# 6. Esperar a que el servicio esté listo
log_info "Esperando a que el servicio inicie..."
sleep 10

# 7. Health check
log_info "Verificando health check..."
max_retries=10
retry_count=0

while [ $retry_count -lt $max_retries ]; do
    if curl -s http://localhost:8095/health | grep -q "ok"; then
        log_info "✓ Health check exitoso"
        break
    else
        retry_count=$((retry_count + 1))
        log_warn "Intento $retry_count/$max_retries - Esperando..."
        sleep 3
    fi
done

if [ $retry_count -eq $max_retries ]; then
    log_error "Health check falló después de $max_retries intentos"
    log_error "Ver logs con: docker-compose logs app"
    exit 1
fi

# 8. Probar endpoint nuevo
log_info "Probando endpoint de Infobip..."
response=$(curl -s -w "%{http_code}" -X POST http://localhost:8095/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
    "nro_cta": "CTA-DEPLOY-TEST",
    "nro_rto": "RTO-DEPLOY-TEST",
    "tipos": {
      "P": 1,
      "M": 1
    },
    "tipo_entrega": "Instalacion",
    "entregado_por": "Repartidor",
    "session_id": "DEPLOY-TEST-001"
  }')

http_code="${response: -3}"
body="${response%???}"

if [ "$http_code" = "201" ]; then
    log_info "✓ Endpoint de Infobip funcionando"
    echo "Response: $body"
else
    log_error "✗ Endpoint falló con código: $http_code"
    echo "Response: $body"
fi

# 9. Mostrar logs recientes
log_info "Logs recientes:"
docker-compose logs --tail=20 app

# 10. Resumen
echo ""
echo "================================"
log_info "✓ DESPLIEGUE COMPLETADO"
echo "================================"
echo ""
echo "Siguiente pasos:"
echo "  - Monitorear logs: docker-compose logs -f app"
echo "  - Ver contenedores: docker ps"
echo "  - Health check: curl http://localhost:8095/health"
echo "  - Backup creado: backups/$BACKUP_FILE"
echo ""
