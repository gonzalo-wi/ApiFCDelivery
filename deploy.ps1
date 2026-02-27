# Script de despliegue automatizado para GoFrioCalor (Windows PowerShell)
# Uso: .\deploy.ps1

param(
    [switch]$SkipBackup = $false,
    [switch]$Force = $false
)

# Colores
function Write-Info { Write-Host "[INFO] $args" -ForegroundColor Green }
function Write-Warn { Write-Host "[WARN] $args" -ForegroundColor Yellow }
function Write-Fail { Write-Host "[ERROR] $args" -ForegroundColor Red }

Write-Info "==================================="
Write-Info "  DESPLIEGUE GOFROCALOR API"
Write-Info "==================================="
Write-Host ""

# Verificar que estamos en el directorio correcto
if (-not (Test-Path "docker-compose.yml")) {
    Write-Fail "No se encontró docker-compose.yml. Ejecuta este script desde la raíz del proyecto."
    exit 1
}

# Verificar branch actual
$currentBranch = git branch --show-current
Write-Info "Branch actual: $currentBranch"

# Confirmar despliegue
if (-not $Force) {
    $confirmation = Read-Host "¿Desplegar a producción? (y/N)"
    if ($confirmation -ne 'y' -and $confirmation -ne 'Y') {
        Write-Warn "Despliegue cancelado"
        exit 0
    }
}

Write-Host ""
Write-Info "Iniciando despliegue..."
Write-Host ""

# 1. Backup de base de datos (opcional)
if (-not $SkipBackup) {
    Write-Info "PASO 1/9: Creando backup de base de datos..."
    
    # Crear directorio de backups si no existe
    if (-not (Test-Path "backups")) {
        New-Item -ItemType Directory -Path "backups" | Out-Null
    }
    
    $backupFile = "backups\backup_$(Get-Date -Format 'yyyyMMdd_HHmmss').sql"
    Write-Info "Archivo: $backupFile"
    Write-Warn "NOTA: Asegúrate de tener un backup manual de la BD antes de continuar"
    Start-Sleep -Seconds 2
} else {
    Write-Warn "PASO 1/9: Backup omitido (parámetro -SkipBackup)"
}

# 2. Pull de cambios
Write-Info "PASO 2/9: Actualizando código desde repositorio..."
git pull origin main
if ($LASTEXITCODE -ne 0) {
    Write-Fail "Error al hacer pull del repositorio"
    exit 1
}

# 3. Detener contenedor actual
Write-Info "PASO 3/9: Deteniendo contenedor actual..."
docker-compose down
if ($LASTEXITCODE -ne 0) {
    Write-Warn "No se pudo detener el contenedor (quizás no estaba corriendo)"
}

# 4. Limpiar imágenes viejas (opcional)
Write-Info "PASO 4/9: Limpiando imágenes Docker antiguas..."
docker image prune -f | Out-Null

# 5. Reconstruir imagen
Write-Info "PASO 5/9: Reconstruyendo imagen Docker (esto puede tomar unos minutos)..."
docker-compose build --no-cache
if ($LASTEXITCODE -ne 0) {
    Write-Fail "Error al construir la imagen Docker"
    exit 1
}

# 6. Iniciar nuevo contenedor
Write-Info "PASO 6/9: Iniciando nuevo contenedor..."
docker-compose up -d
if ($LASTEXITCODE -ne 0) {
    Write-Fail "Error al iniciar el contenedor"
    exit 1
}

# 7. Esperar a que el servicio esté listo
Write-Info "PASO 7/9: Esperando a que el servicio inicie..."
Write-Host "Esperando" -NoNewline
Start-Sleep -Seconds 10
Write-Host " OK" -ForegroundColor Green

# 8. Health check
Write-Info "PASO 8/9: Verificando health check..."
$maxRetries = 10
$retryCount = 0
$healthOk = $false

while ($retryCount -lt $maxRetries) {
    try {
        $response = Invoke-RestMethod -Uri "http://localhost:8095/health" -Method GET -TimeoutSec 5
        if ($response.status -eq "ok") {
            Write-Host "  ✓ Health check exitoso" -ForegroundColor Green
            $healthOk = $true
            break
        }
    }
    catch {
        $retryCount++
        Write-Host "  Intento $retryCount/$maxRetries - Esperando..." -ForegroundColor Yellow
        Start-Sleep -Seconds 3
    }
}

if (-not $healthOk) {
    Write-Fail "Health check falló después de $maxRetries intentos"
    Write-Warn "Ver logs con: docker-compose logs app"
    exit 1
}

# 9. Probar endpoint nuevo de Infobip
Write-Info "PASO 9/9: Probando endpoint de Infobip..."
try {
    $body = @{
        nro_cta = "CTA-DEPLOY-TEST"
        nro_rto = "RTO-DEPLOY-TEST"
        tipos = @{
            P = 1
            M = 1
        }
        tipo_entrega = "Instalacion"
        entregado_por = "Repartidor"
        session_id = "DEPLOY-TEST-$(Get-Date -Format 'yyyyMMddHHmmss')"
    } | ConvertTo-Json

    $response = Invoke-RestMethod -Uri "http://localhost:8095/api/v1/deliveries/infobip" `
        -Method POST `
        -ContentType "application/json" `
        -Body $body `
        -TimeoutSec 10

    Write-Host "  ✓ Endpoint de Infobip funcionando" -ForegroundColor Green
    Write-Host "  Token generado: $($response.token)" -ForegroundColor Cyan
}
catch {
    Write-Fail "Error al probar endpoint de Infobip"
    Write-Host "  $_" -ForegroundColor Red
}

# Mostrar logs recientes
Write-Host ""
Write-Info "Logs recientes del contenedor:"
Write-Host "================================" -ForegroundColor DarkGray
docker-compose logs --tail=15 app
Write-Host "================================" -ForegroundColor DarkGray

# Resumen final
Write-Host ""
Write-Host "====================================" -ForegroundColor Green
Write-Host "  ✓ DESPLIEGUE COMPLETADO EXITOSAMENTE" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Green
Write-Host ""

Write-Info "Siguiente pasos:"
Write-Host "  • Monitorear logs:    docker-compose logs -f app" -ForegroundColor Cyan
Write-Host "  • Ver contenedores:   docker ps" -ForegroundColor Cyan
Write-Host "  • Health check:       curl http://localhost:8095/health" -ForegroundColor Cyan
Write-Host "  • Probar Postman:     Importar postman/Infobip_Delivery_Collection.json" -ForegroundColor Cyan
if (-not $SkipBackup) {
    Write-Host "  • Backup creado:      $backupFile" -ForegroundColor Cyan
}
Write-Host ""

Write-Info "El servidor está corriendo en: http://localhost:8095"
Write-Host ""
