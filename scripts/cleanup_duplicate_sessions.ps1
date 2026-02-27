# Script para limpiar session_ids duplicados antes de aplicar índice único
# Fecha: 2026-02-27

Write-Host "=== Limpieza de Session IDs Duplicados ===" -ForegroundColor Cyan

# Para propósitos de testing, usar valores por defecto
$dbHost = "localhost"
$dbPort = "3306"
$dbUser = "root"
$dbPassword = ""
$dbName = "gofrocalor"

$dbConnStr = "${dbHost}:${dbPort}/${dbName}"
Write-Host "Base de datos: $dbConnStr" -ForegroundColor Yellow

# Opción simple: Setear a NULL los session_id duplicados (conservando el más antiguo)
Write-Host "`nLimpiando session_ids duplicados (conservando el registro más antiguo)..." -ForegroundColor Yellow

$updateQuery = "UPDATE deliveries d2 INNER JOIN (SELECT session_id, MIN(id) as min_id FROM deliveries WHERE session_id IS NOT NULL AND session_id != '' GROUP BY session_id HAVING COUNT(*) > 1) d1 ON d2.session_id = d1.session_id AND d2.id > d1.min_id SET d2.session_id = NULL;"

$result = echo $updateQuery | mysql -h $dbHost -P $dbPort -u $dbUser $dbName 2>&1

if ($LASTEXITCODE -eq 0) {
    Write-Host "Duplicados limpiados exitosamente" -ForegroundColor Green
} else {
    Write-Host "Error o no hay duplicados: $result" -ForegroundColor Yellow
}

Write-Host "`nVerificando que no queden duplicados..." -ForegroundColor Yellow
$checkQuery = "SELECT session_id, COUNT(*) as count FROM deliveries WHERE session_id IS NOT NULL AND session_id != '' GROUP BY session_id HAVING COUNT(*) > 1;"
$check = echo $checkQuery | mysql -h $dbHost -P $dbPort -u $dbUser $dbName 2>&1

if ($check -match "session_id") {
    Write-Host "Aún hay duplicados:" -ForegroundColor Red
    echo $check
} else {
    Write-Host "No hay duplicados. OK para aplicar migración 004" -ForegroundColor Green
}
