# Script para verificar optimizaciones de performance
# Ejecutar: .\tests\test_optimizations.ps1

Write-Host "`n=== VERIFICACIÓN DE OPTIMIZACIONES ===" -ForegroundColor Cyan

# 1. Verificar que el servidor esté corriendo
Write-Host "`n[1] Verificando servidor..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8095/health" -Method Get -ErrorAction Stop
    Write-Host "✓ Servidor activo" -ForegroundColor Green
} catch {
    Write-Host "✗ Servidor no responde. Iniciar con: go run api/cmd/main.go" -ForegroundColor Red
    exit 1
}

# 2. Verificar índices en MySQL
Write-Host "`n[2] Verificando índices de base de datos..." -ForegroundColor Yellow
if (Test-Path ".env") {
    Get-Content .env | ForEach-Object {
        if ($_ -match '^([^=]+)=(.*)$') {
            $key = $matches[1]
            $value = $matches[2]
            [Environment]::SetEnvironmentVariable($key, $value, "Process")
        }
    }
}

$DB_HOST = $env:DB_HOST
$DB_PORT = $env:DB_PORT
$DB_USER = $env:DB_USER
$DB_PASSWORD = $env:DB_PASSWORD
$DB_NAME = $env:DB_NAME

$query = "SHOW INDEX FROM deliveries WHERE Key_name LIKE 'idx_deliveries_%';"
$command = "mysql -h $DB_HOST -P $DB_PORT -u $DB_USER"
if ($DB_PASSWORD) {
    $command += " -p$DB_PASSWORD"
}
$command += " $DB_NAME -e `"$query`""

try {
    $indexes = Invoke-Expression $command 2>&1
    if ($indexes -match "idx_deliveries_token_validation") {
        Write-Host "✓ Índice idx_deliveries_token_validation encontrado" -ForegroundColor Green
    } else {
        Write-Host "✗ Índice idx_deliveries_token_validation NO encontrado" -ForegroundColor Red
        Write-Host "  Ejecutar: .\scripts\apply_migration_006.ps1" -ForegroundColor Yellow
    }
    if ($indexes -match "idx_deliveries_estado") {
        Write-Host "✓ Índice idx_deliveries_estado encontrado" -ForegroundColor Green
    }
    if ($indexes -match "idx_deliveries_fecha_cuenta") {
        Write-Host "✓ Índice idx_deliveries_fecha_cuenta encontrado" -ForegroundColor Green
    }
    if ($indexes -match "idx_deliveries_fecha_rto") {
        Write-Host "✓ Índice idx_deliveries_fecha_rto encontrado" -ForegroundColor Green
    }
} catch {
    Write-Host "✗ No se pudo verificar índices (mysql no disponible)" -ForegroundColor Yellow
}

# 3. Test de ValidateToken (debe usar FindByTokenAndFilters)
Write-Host "`n[3] Probando ValidateToken optimizado..." -ForegroundColor Yellow

$validateRequest = @{
    token = "test-token-123"
    nroCta = "12345"
    fechaAccion = "2026-03-03"
} | ConvertTo-Json

Write-Host "Request: POST /api/v1/mobile/validate-token" -ForegroundColor White
$startTime = Get-Date
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8095/api/v1/mobile/validate-token" `
        -Method Post `
        -Body $validateRequest `
        -ContentType "application/json" `
        -ErrorAction Stop
    
    $endTime = Get-Date
    $elapsed = ($endTime - $startTime).TotalMilliseconds
    
    Write-Host "✓ ValidateToken respondió en ${elapsed}ms" -ForegroundColor Green
    if ($elapsed -lt 50) {
        Write-Host "  Excelente performance (<50ms)" -ForegroundColor Green
    } elseif ($elapsed -lt 100) {
        Write-Host "  Buena performance (<100ms)" -ForegroundColor Cyan
    } else {
        Write-Host "  Performance normal (>100ms) - Verificar índices" -ForegroundColor Yellow
    }
    
    if ($response.Valid) {
        Write-Host "  Token válido encontrado" -ForegroundColor White
        Write-Host "  NroRto: $($response.NroRto)" -ForegroundColor White
        Write-Host "  Cantidad: $($response.Cantidad)" -ForegroundColor White
    } else {
        Write-Host "  Token no encontrado (esperado si no existe en DB)" -ForegroundColor White
    }
} catch {
    Write-Host "✗ Error en ValidateToken: $($_.Exception.Message)" -ForegroundColor Red
}

# 4. Test de CompleteDelivery (debe usar HashMap)
Write-Host "`n[4] Probando CompleteDelivery optimizado..." -ForegroundColor Yellow

$completeRequest = @{
    token = "test-token-123"
    validated = @("DISPENSER-001", "DISPENSER-002")
    observacion = "Test de optimización"
} | ConvertTo-Json

Write-Host "Request: POST /api/v1/mobile/complete-delivery" -ForegroundColor White
$startTime = Get-Date
try {
    $response = Invoke-RestMethod -Uri "http://localhost:8095/api/v1/mobile/complete-delivery" `
        -Method Post `
        -Body $completeRequest `
        -ContentType "application/json" `
        -ErrorAction Stop
    
    $endTime = Get-Date
    $elapsed = ($endTime - $startTime).TotalMilliseconds
    
    Write-Host "✓ CompleteDelivery respondió en ${elapsed}ms" -ForegroundColor Green
    if ($elapsed -lt 100) {
        Write-Host "  Excelente performance (<100ms)" -ForegroundColor Green
    } elseif ($elapsed -lt 200) {
        Write-Host "  Buena performance (<200ms)" -ForegroundColor Cyan
    } else {
        Write-Host "  Performance normal (>200ms)" -ForegroundColor Yellow
    }
    
    # Verificar campo workOrderQueued
    if ($null -ne $response.WorkOrderQueued) {
        Write-Host "✓ Campo workOrderQueued presente: $($response.WorkOrderQueued)" -ForegroundColor Green
    } else {
        Write-Host "✗ Campo workOrderQueued faltante (código no actualizado)" -ForegroundColor Red
    }
    
    Write-Host "  Estado: $($response.Estado)" -ForegroundColor White
    Write-Host "  Entregados: $($response.DispensersEntregados)" -ForegroundColor White
} catch {
    $errorMessage = $_.ErrorDetails.Message
    if ($errorMessage -match "No existe una entrega") {
        Write-Host "  Token no encontrado en DB (esperado si no existe)" -ForegroundColor White
    } else {
        Write-Host "✗ Error: $errorMessage" -ForegroundColor Red
    }
}

# 5. Resumen
Write-Host "`n=== RESUMEN ===" -ForegroundColor Cyan
Write-Host "Para verificar optimizaciones completas:" -ForegroundColor White
Write-Host "  1. Los índices deben estar aplicados (ejecutar .\scripts\apply_migration_006.ps1)" -ForegroundColor White
Write-Host "  2. ValidateToken debe responder en <50ms con índices" -ForegroundColor White
Write-Host "  3. CompleteDelivery debe incluir workOrderQueued en respuesta" -ForegroundColor White
Write-Host "  4. El servidor debe usar FindByTokenAndFilters (revisar logs)" -ForegroundColor White
Write-Host "`nVer logs del servidor para más detalles sobre queries ejecutadas." -ForegroundColor Yellow
