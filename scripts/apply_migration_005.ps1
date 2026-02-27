# Script para aplicar la migración 005: Agregar información del cliente a deliveries
# Fecha: 2026-02-27

Write-Host "=== Aplicando Migración 005: Agregar Name, Address, Locality a Deliveries ===" -ForegroundColor Cyan

# Cargar variables de entorno
if (Test-Path .env) {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            $key = $matches[1].Trim()
            $value = $matches[2].Trim()
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
    Write-Host "Variables de entorno cargadas desde .env" -ForegroundColor Green
}

# Obtener variables de conexión
$dbUser = $env:DB_USER
$dbPassword = $env:DB_PASSWORD
$dbHost = $env:DB_HOST
$dbPort = $env:DB_PORT
$dbName = $env:DB_NAME

if (-not $dbUser -or -not $dbPassword -or -not $dbHost -or -not $dbPort -or -not $dbName) {
    Write-Host "Error: Variables de entorno de base de datos no encontradas" -ForegroundColor Red
    exit 1
}

Write-Host "Conectando a base de datos: ${dbHost}:${dbPort}/${dbName}" -ForegroundColor Yellow

# Leer el archivo SQL de migración
$migrationFile = "migrations/005_add_client_info_to_deliveries.sql"
if (-not (Test-Path $migrationFile)) {
    Write-Host "Error: Archivo de migración no encontrado: $migrationFile" -ForegroundColor Red
    exit 1
}

$sqlContent = Get-Content $migrationFile -Raw
Write-Host "Archivo de migración cargado" -ForegroundColor Green

# Ejecutar migración
Write-Host ""
Write-Host "Ejecutando migración..." -ForegroundColor Yellow
try {
    $sqlContent | & mysql -h $dbHost -P $dbPort -u $dbUser -p"$dbPassword" $dbName
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host ""
        Write-Host "Migración 005 aplicada exitosamente!" -ForegroundColor Green
        Write-Host ""
        Write-Host "Cambios realizados:" -ForegroundColor Cyan
        Write-Host "  - Columna 'name' agregada (VARCHAR 200)" -ForegroundColor White
        Write-Host "  - Columna 'address' agregada (VARCHAR 300)" -ForegroundColor White
        Write-Host "  - Columna 'locality' agregada (VARCHAR 100)" -ForegroundColor White
        Write-Host "  - Índice 'idx_deliveries_name' creado" -ForegroundColor White
        Write-Host "  - Índice 'idx_deliveries_locality' creado" -ForegroundColor White
        Write-Host ""
        
        # Verificar las columnas
        Write-Host "Verificando estructura de tabla deliveries..." -ForegroundColor Yellow
        "DESCRIBE deliveries;" | & mysql -h $dbHost -P $dbPort -u $dbUser -p"$dbPassword" $dbName
        
    } else {
        Write-Host ""
        Write-Host "Error al aplicar la migración" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host ""
    Write-Host "Error: $_" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Proceso completado!" -ForegroundColor Green
