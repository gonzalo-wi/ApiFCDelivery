# Script para aplicar migración 007: Audit Events Table
# Ejecutar: .\scripts\apply_migration_007.ps1

Write-Host "Aplicando Migration 007: Audit Events Table" -ForegroundColor Cyan

# Cargar variables de entorno desde .env
if (Test-Path ".env") {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            $key = $matches[1]
            $value = $matches[2]
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "Variables de entorno cargadas desde .env" -ForegroundColor Green
} else {
    Write-Host "Archivo .env no encontrado, usando variables del sistema" -ForegroundColor Yellow
}

$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_USER = $env:DB_USER
$DB_PASSWORD = $env:DB_PASSWORD
$DB_NAME = $env:DB_NAME

Write-Host "Conectando a PostgreSQL: ${DB_USER}@${DB_HOST}:${DB_PORT}/${DB_NAME}" -ForegroundColor Yellow

# Construir connection string para PostgreSQL
$env:PGPASSWORD = $DB_PASSWORD
$psqlCommand = "psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f migrations/007_create_audit_events.sql"

Write-Host "Ejecutando migración..." -ForegroundColor Cyan

try {
    # Ejecutar la migración
    Invoke-Expression $psqlCommand
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`nMigración 007 aplicada exitosamente!" -ForegroundColor Green
        Write-Host "Tabla 'audit_events' creada con:" -ForegroundColor Green
        Write-Host "  - 13 columnas (id, occurred_at, service, entity_type, entity_id, action, actor_type, actor_id, request_id, trace_id, ip_address, user_agent, before_state, after_state, metadata)" -ForegroundColor Cyan
        Write-Host "  - 7 índices optimizados para consultas frecuentes" -ForegroundColor Cyan
        Write-Host "  - Soporte para particionamiento por fecha" -ForegroundColor Cyan
        Write-Host "  - Función cleanup_old_audit_events() para retención de datos" -ForegroundColor Cyan
    } else {
        Write-Host "`nError al aplicar la migración" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "`nError: $_" -ForegroundColor Red
    exit 1
} finally {
    # Limpiar la variable de password
    Remove-Item Env:\PGPASSWORD -ErrorAction SilentlyContinue
}

Write-Host "`nPróximos pasos:" -ForegroundColor Yellow
Write-Host "1. Iniciar el servidor: .\start_server.ps1" -ForegroundColor White
Write-Host "2. Probar endpoints de auditoría:" -ForegroundColor White
Write-Host "   GET  /dispenser-operations/api/v1/audit/recent" -ForegroundColor Cyan
Write-Host "   GET  /dispenser-operations/api/v1/audit/entity/:type/:id" -ForegroundColor Cyan
Write-Host "   POST /dispenser-operations/api/v1/audit/search" -ForegroundColor Cyan
Write-Host "   GET  /dispenser-operations/api/v1/audit/stats" -ForegroundColor Cyan
Write-Host "3. Integrar logging en handlers existentes (delivery, work_order, etc.)" -ForegroundColor White
