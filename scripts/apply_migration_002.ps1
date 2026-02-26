# Script para aplicar la migración de session_id a deliveries
# Ejecutar desde la raíz del proyecto

$env:MYSQL_HOST = "localhost"
$env:MYSQL_PORT = "3306"
$env:MYSQL_USER = "root"
$env:MYSQL_PASSWORD = "tu_password"
$env:MYSQL_DATABASE = "gofriocalor"

Write-Host "Aplicando migración 002_add_session_id_to_deliveries.sql..." -ForegroundColor Cyan

# Ejecutar la migración
Get-Content ".\migrations\002_add_session_id_to_deliveries.sql" | mysql -h $env:MYSQL_HOST -P $env:MYSQL_PORT -u $env:MYSQL_USER -p$env:MYSQL_PASSWORD $env:MYSQL_DATABASE

if ($LASTEXITCODE -eq 0) {
    Write-Host "✅ Migración aplicada exitosamente" -ForegroundColor Green
} else {
    Write-Host "❌ Error al aplicar la migración" -ForegroundColor Red
}
