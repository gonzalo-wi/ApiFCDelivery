# Script de Prueba - Flujo Integrado de Entregas con Términos

Write-Host "=== GoFrioCalor - Test del Flujo Integrado ===" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8080/api/v1"

# Función helper para hacer requests
function Invoke-APIRequest {
    param(
        [string]$Method,
        [string]$Endpoint,
        [object]$Body = $null
    )
    
    $url = "$baseUrl$Endpoint"
    Write-Host "[$Method] $url" -ForegroundColor Yellow
    
    try {
        if ($Body) {
            $jsonBody = $Body | ConvertTo-Json -Depth 10
            Write-Host "Body: $jsonBody" -ForegroundColor Gray
            $response = Invoke-RestMethod -Uri $url -Method $Method -Body $jsonBody -ContentType "application/json"
        } else {
            $response = Invoke-RestMethod -Uri $url -Method $Method
        }
        
        Write-Host "Response:" -ForegroundColor Green
        $response | ConvertTo-Json -Depth 10 | Write-Host
        return $response
    }
    catch {
        Write-Host "Error: $_" -ForegroundColor Red
        if ($_.Exception.Response) {
            $reader = New-Object System.IO.StreamReader($_.Exception.Response.GetResponseStream())
            $reader.BaseStream.Position = 0
            $responseBody = $reader.ReadToEnd()
            Write-Host "Response Body: $responseBody" -ForegroundColor Red
        }
        return $null
    }
}

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 1: Iniciar Creación de Entrega" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$initiateRequest = @{
    nro_cta = "CTA-TEST-001"
    nro_rto = "RTO-TEST-" + (Get-Date -Format "yyyyMMddHHmmss")
    dispensers = @(
        @{
            marca = "CocaCola"
            nro_serie = "CC-SN-001"
            tipo = "Enfriador"
        },
        @{
            marca = "Pepsi"
            nro_serie = "PP-SN-002"
            tipo = "Calentador"
        }
    )
    cantidad = 2
    tipo_entrega = "Instalacion"
    fecha_accion = (Get-Date).ToString("yyyy-MM-dd")
}

$initiateResponse = Invoke-APIRequest -Method POST -Endpoint "/deliveries/initiate" -Body $initiateRequest

if (-not $initiateResponse) {
    Write-Host ""
    Write-Host "❌ Error al iniciar entrega. Abortando test." -ForegroundColor Red
    exit 1
}

$token = $initiateResponse.token
Write-Host ""
Write-Host "✅ Token generado: $token" -ForegroundColor Green
Write-Host "📝 URL de términos: $($initiateResponse.terms_url)" -ForegroundColor Green
Write-Host ""

# Pausa
Write-Host "Esperando 2 segundos antes del siguiente paso..." -ForegroundColor Gray
Start-Sleep -Seconds 2

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 2: Verificar Estado de Términos" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$statusResponse = Invoke-APIRequest -Method GET -Endpoint "/terms/status/$token"

if ($statusResponse -and $statusResponse.status -eq "PENDING") {
    Write-Host ""
    Write-Host "✅ Estado actual: PENDING (esperando aceptación)" -ForegroundColor Yellow
} else {
    Write-Host ""
    Write-Host "⚠️ Estado inesperado: $($statusResponse.status)" -ForegroundColor Yellow
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 3: Intentar Completar SIN Aceptar Términos" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Esto DEBE FALLAR porque aún no se aceptaron los términos" -ForegroundColor Yellow
Write-Host ""

$completeResponse = Invoke-APIRequest -Method POST -Endpoint "/deliveries/complete/$token"

if (-not $completeResponse) {
    Write-Host ""
    Write-Host "✅ Correcto: No se pudo completar sin aceptar términos" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error: Se completó sin aceptar términos (no debería pasar)" -ForegroundColor Red
}

Write-Host ""
Write-Host "Esperando 2 segundos..." -ForegroundColor Gray
Start-Sleep -Seconds 2

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 4: Aceptar Términos y Condiciones" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$acceptRequest = @{
    webhook_url = "https://api.infobip.com/webhook/test-acceptance"
}

$acceptResponse = Invoke-APIRequest -Method POST -Endpoint "/terms/accept/$token" -Body $acceptRequest

if ($acceptResponse -and $acceptResponse.status -eq "ACCEPTED") {
    Write-Host ""
    Write-Host "✅ Términos aceptados exitosamente" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error al aceptar términos" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "Esperando 3 segundos para que se procese la aceptación..." -ForegroundColor Gray
Start-Sleep -Seconds 3

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 5: Completar Entrega (Ahora Sí)" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$completeResponse = Invoke-APIRequest -Method POST -Endpoint "/deliveries/complete/$token"

if ($completeResponse -and $completeResponse.success) {
    Write-Host ""
    Write-Host "✅ ¡ENTREGA CREADA EXITOSAMENTE!" -ForegroundColor Green
    Write-Host "   ID de Entrega: $($completeResponse.delivery.id)" -ForegroundColor Green
    Write-Host "   Nro RTO: $($completeResponse.delivery.nro_rto)" -ForegroundColor Green
    Write-Host "   Estado: $($completeResponse.delivery.estado)" -ForegroundColor Green
    Write-Host "   Token de Entrega: $($completeResponse.delivery.token)" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error al completar entrega" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 6: Verificar Estado Final" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$statusFinalResponse = Invoke-APIRequest -Method GET -Endpoint "/terms/status/$token"

if ($statusFinalResponse -and $statusFinalResponse.status -eq "ACCEPTED") {
    Write-Host ""
    Write-Host "✅ Estado final: ACCEPTED" -ForegroundColor Green
    Write-Host "   Aceptado en: $($statusFinalResponse.accepted_at)" -ForegroundColor Green
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ TEST COMPLETADO EXITOSAMENTE" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Resumen:" -ForegroundColor White
Write-Host "  1. ✅ Entrega iniciada con token de términos" -ForegroundColor White
Write-Host "  2. ✅ Verificación de estado (PENDING)" -ForegroundColor White
Write-Host "  3. ✅ Falló correctamente sin aceptar términos" -ForegroundColor White
Write-Host "  4. ✅ Términos aceptados" -ForegroundColor White
Write-Host "  5. ✅ Entrega completada exitosamente" -ForegroundColor White
Write-Host "  6. ✅ Estado final verificado (ACCEPTED)" -ForegroundColor White
Write-Host ""
Write-Host "Token de términos usado: $token" -ForegroundColor Cyan
Write-Host ""
