# Script para aplicar migración 008 - Agregar campos de email
# Agrega columnas email a deliveries y work_orders

$ErrorActionPreference = "Stop"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Aplicando Migración 008: Email Fields" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Cargar variables de entorno
if (Test-Path ".env") {
    Get-Content ".env" | ForEach-Object {
        if ($_ -match "^\s*([^#][^=]+)\s*=\s*(.+)\s*$") {
            $name = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($name, $value, "Process")
            Write-Host "✓ Loaded $name" -ForegroundColor Green
        }
    }
} else {
    Write-Host "⚠️  Archivo .env no encontrado" -ForegroundColor Yellow
    exit 1
}

# Obtener credenciales de base de datos
$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_USER = $env:DB_USER
$DB_PASSWORD = $env:DB_PASSWORD
$DB_NAME = $env:DB_NAME

if (-not $DB_HOST -or -not $DB_PORT -or -not $DB_USER -or -not $DB_PASSWORD -or -not $DB_NAME) {
    Write-Host "❌ Variables de base de datos incompletas" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Conectando a: $DB_HOST:$DB_PORT/$DB_NAME" -ForegroundColor Yellow
Write-Host ""

# Construir connection string para psql
$env:PGPASSWORD = $DB_PASSWORD

# Verificar si psql está disponible
$psqlPath = Get-Command psql -ErrorAction SilentlyContinue
if (-not $psqlPath) {
    Write-Host "❌ psql no está instalado o no está en el PATH" -ForegroundColor Red
    Write-Host ""
    Write-Host "Alternativas:" -ForegroundColor Yellow
    Write-Host "1. Instalar PostgreSQL client tools" -ForegroundColor White
    Write-Host "2. Usar script Go: go run scripts/run_migration_008.go" -ForegroundColor White
    exit 1
}

try {
    # Ejecutar migración
    Write-Host "Aplicando migración..." -ForegroundColor Cyan
    
    $migrationFile = "migrations/008_add_email_fields.sql"
    
    if (-not (Test-Path $migrationFile)) {
        Write-Host "❌ Archivo de migración no encontrado: $migrationFile" -ForegroundColor Red
        exit 1
    }
    
    psql -h $DB_HOST -p $DB_PORT -U $DB_USER -d $DB_NAME -f $migrationFile
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "✅ Migración aplicada exitosamente!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Cambios realizados:" -ForegroundColor Cyan
        Write-Host "  • Agregada columna 'email' a tabla 'deliveries'" -ForegroundColor White
        Write-Host "  • Agregada columna 'email' a tabla 'work_orders'" -ForegroundColor White
        Write-Host "  • Creados índices para búsquedas por email" -ForegroundColor White
        Write-Host ""
    } else {
        Write-Host ""
        Write-Host "❌ Error al aplicar migración" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host ""
    Write-Host "❌ Error: $_" -ForegroundColor Red
    exit 1
} finally {
    # Limpiar variable de password
    $env:PGPASSWORD = $null
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Migración completada" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
