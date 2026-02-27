# Script para aplicar la migraci\u00f3n 004: \u00cdndice \u00fanico en session_id
# Fecha: 2026-02-27

Write-Host "=== Aplicando Migraci\u00f3n 004: \u00cdndice \u00danico en session_id ===" -ForegroundColor Cyan

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

# Obtener variables de conexi\u00f3n
$dbUser = $env:DB_USER
$dbPassword = $env:DB_PASSWORD
$dbHost = $env:DB_HOST
$dbPort = $env:DB_PORT
$dbName = $env:DB_NAME

if (-not $dbUser -or -not $dbPassword -or -not $dbHost -or -not $dbPort -or -not $dbName) {
    Write-Host "Error: Variables de entorno de base de datos no encontradas" -ForegroundColor Red
    exit 1
}

Write-Host "Conectando a base de datos: $dbHost:$dbPort/$dbName" -ForegroundColor Yellow

# Construir comando mysql
$migrationFile = "migrations/004_add_unique_session_id.sql"

if (-not (Test-Path $migrationFile)) {
    Write-Host "Error: Archivo de migraci\u00f3n no encontrado: $migrationFile" -ForegroundColor Red
    exit 1
}

Write-Host "Aplicando migraci\u00f3n: $migrationFile" -ForegroundColor Yellow

# Ejecutar migraci\u00f3n
$mysqlCmd = "mysql -h$dbHost -P$dbPort -u$dbUser -p$dbPassword $dbName"

try {
    Get-Content $migrationFile | & mysql -h$dbHost -P$dbPort -u$dbUser "-p$dbPassword" $dbName
    
    if ($LASTEXITCODE -eq 0) {
        Write-Host "`n\u2713 Migraci\u00f3n aplicada exitosamente" -ForegroundColor Green
        
        # Verificar el \u00edndice
        Write-Host "`nVerificando \u00edndice creado..." -ForegroundColor Yellow
        $verifyQuery = "SHOW INDEX FROM deliveries WHERE Column_name = 'session_id';"
        echo $verifyQuery | & mysql -h$dbHost -P$dbPort -u$dbUser "-p$dbPassword" $dbName
        
        Write-Host "`n=== Migraci\u00f3n Completada ===" -ForegroundColor Green
        Write-Host "Ahora el campo session_id tiene un \u00edndice \u00fanico que previene duplicados" -ForegroundColor Cyan
    } else {
        Write-Host "`nError aplicando migraci\u00f3n (Exit Code: $LASTEXITCODE)" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "`nError: $_" -ForegroundColor Red
    exit 1
}
