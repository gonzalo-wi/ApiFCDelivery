# Script de prueba para el endpoint de Infobip Delivery
# Ejecuta: .\test_infobip_delivery.ps1

$baseUrl = "http://localhost:8080/api/v1"
$infobipEndpoint = "$baseUrl/deliveries/infobip"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  PRUEBAS DE ENDPOINT INFOBIP DELIVERY  " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Función para hacer requests
function Invoke-TestRequest {
    param(
        [string]$TestName,
        [hashtable]$Body,
        [int]$ExpectedStatus = 201
    )
    
    Write-Host "TEST: $TestName" -ForegroundColor Yellow
    Write-Host "Request Body:" -ForegroundColor Gray
    Write-Host ($Body | ConvertTo-Json) -ForegroundColor Gray
    
    try {
        $response = Invoke-WebRequest -Uri $infobipEndpoint `
            -Method POST `
            -ContentType "application/json" `
            -Body ($Body | ConvertTo-Json) `
            -UseBasicParsing
        
        if ($response.StatusCode -eq $ExpectedStatus) {
            Write-Host "✓ PASS" -ForegroundColor Green
            $content = $response.Content | ConvertFrom-Json
            Write-Host "Response:" -ForegroundColor Green
            Write-Host ($content | ConvertTo-Json) -ForegroundColor Green
            Write-Host ""
            return $true
        } else {
            Write-Host "✗ FAIL - Expected $ExpectedStatus, got $($response.StatusCode)" -ForegroundColor Red
            Write-Host ""
            return $false
        }
    }
    catch {
        $statusCode = $_.Exception.Response.StatusCode.value__
        if ($statusCode -eq $ExpectedStatus) {
            Write-Host "✓ PASS (Expected error)" -ForegroundColor Green
            Write-Host "Error Response:" -ForegroundColor Yellow
            Write-Host $_.Exception.Message -ForegroundColor Yellow
            Write-Host ""
            return $true
        } else {
            Write-Host "✗ FAIL - Unexpected error" -ForegroundColor Red
            Write-Host $_.Exception.Message -ForegroundColor Red
            Write-Host ""
            return $false
        }
    }
}

# Verificar que el servidor está corriendo
Write-Host "Verificando conexión al servidor..." -ForegroundColor Cyan
try {
    $healthCheck = Invoke-WebRequest -Uri "http://localhost:8080/health" -UseBasicParsing
    Write-Host "✓ Servidor activo" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "✗ Error: El servidor no está corriendo en http://localhost:8080" -ForegroundColor Red
    Write-Host "Inicie el servidor primero con: go run api/cmd/main.go" -ForegroundColor Yellow
    exit 1
}

# Contador de tests
$totalTests = 0
$passedTests = 0

# TEST 1: Crear entrega con 2 dispensers de pie
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Instalación con 2 dispensers de pie" -Body @{
    nro_cta = "CTA12345"
    nro_rto = "RTO001"
    tipos = @{
        P = 2
        M = 0
    }
    tipo_entrega = "Instalacion"
    entregado_por = "Repartidor"
    session_id = "INF-TEST-001"
    fecha_accion = "2026-02-25"
}
if ($result) { $passedTests++ }

# TEST 2: Crear entrega con 1 de cada tipo
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Recambio con 1 de pie y 1 de mesada" -Body @{
    nro_cta = "CTA99999"
    nro_rto = "RTO999"
    tipos = @{
        P = 1
        M = 1
    }
    tipo_entrega = "Recambio"
    entregado_por = "Tecnico"
    session_id = "INF-TEST-002"
}
if ($result) { $passedTests++ }

# TEST 3: Crear entrega solo con dispensers de mesada
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Retiro con 3 dispensers de mesada" -Body @{
    nro_cta = "CTA77777"
    nro_rto = "RTO777"
    tipos = @{
        P = 0
        M = 3
    }
    tipo_entrega = "Retiro"
    entregado_por = "Repartidor"
    session_id = "INF-TEST-003"
    fecha_accion = "2026-02-26T10:30:00Z"
}
if ($result) { $passedTests++ }

# TEST 4: Error - Sin dispensers
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Error: Sin dispensers especificados" -Body @{
    nro_cta = "CTA00000"
    nro_rto = "RTO000"
    tipos = @{
        P = 0
        M = 0
    }
    tipo_entrega = "Instalacion"
    entregado_por = "Repartidor"
    session_id = "INF-TEST-004"
} -ExpectedStatus 400
if ($result) { $passedTests++ }

# TEST 5: Error - Falta campo requerido
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Error: Falta nro_cta" -Body @{
    nro_rto = "RTO555"
    tipos = @{
        P = 1
        M = 0
    }
    tipo_entrega = "Instalacion"
    entregado_por = "Repartidor"
    session_id = "INF-TEST-005"
} -ExpectedStatus 400
if ($result) { $passedTests++ }

# TEST 6: Error - Tipo de entrega inválido
Write-Host "----------------------------------------" -ForegroundColor DarkGray
$totalTests++
$result = Invoke-TestRequest -TestName "Error: Tipo de entrega inválido" -Body @{
    nro_cta = "CTA11111"
    nro_rto = "RTO111"
    tipos = @{
        P = 1
        M = 1
    }
    tipo_entrega = "TipoInvalido"
    entregado_por = "Repartidor"
    session_id = "INF-TEST-006"
} -ExpectedStatus 400
if ($result) { $passedTests++ }

# TEST 7: Múltiples requests simultáneos (concurrencia)
Write-Host "----------------------------------------" -ForegroundColor DarkGray
Write-Host "TEST: Concurrencia - 5 requests simultáneos" -ForegroundColor Yellow

$jobs = @()
for ($i = 1; $i -le 5; $i++) {
    $body = @{
        nro_cta = "CTA-CONCURRENT-$i"
        nro_rto = "RTO-CONCURRENT-$i"
        tipos = @{
            P = 1
            M = 1
        }
        tipo_entrega = "Instalacion"
        entregado_por = "Repartidor"
        session_id = "INF-CONCURRENT-$i"
    } | ConvertTo-Json
    
    $job = Start-Job -ScriptBlock {
        param($url, $body)
        Invoke-WebRequest -Uri $url -Method POST -ContentType "application/json" -Body $body -UseBasicParsing
    } -ArgumentList $infobipEndpoint, $body
    
    $jobs += $job
}

Write-Host "Esperando respuestas..." -ForegroundColor Gray
$results = $jobs | Wait-Job | Receive-Job
$jobs | Remove-Job

$successCount = ($results | Where-Object { $_.StatusCode -eq 201 }).Count
$totalTests++

if ($successCount -eq 5) {
    Write-Host "✓ PASS - Todas las requests concurrentes exitosas ($successCount/5)" -ForegroundColor Green
    $passedTests++
} else {
    Write-Host "✗ FAIL - Solo $successCount/5 requests exitosas" -ForegroundColor Red
}
Write-Host ""

# Resumen final
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "           RESUMEN DE PRUEBAS           " -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Total de tests: $totalTests" -ForegroundColor White
Write-Host "Tests exitosos: $passedTests" -ForegroundColor Green
Write-Host "Tests fallidos: $($totalTests - $passedTests)" -ForegroundColor $(if ($totalTests -eq $passedTests) { "Green" } else { "Red" })
Write-Host ""

if ($totalTests -eq $passedTests) {
    Write-Host "✓ TODOS LOS TESTS PASARON" -ForegroundColor Green
    exit 0
} else {
    Write-Host "✗ ALGUNOS TESTS FALLARON" -ForegroundColor Red
    exit 1
}
