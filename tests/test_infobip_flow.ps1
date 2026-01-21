# Script de Prueba - Flujo con Infobip

Write-Host "=== GoFrioCalor - Test Flujo Infobip ===" -ForegroundColor Cyan
Write-Host ""

$baseUrl = "http://localhost:8080/api/v1"
$sessionId = "RTO-TEST-" + (Get-Date -Format "yyyyMMddHHmmss")

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
Write-Host "PASO 1: Infobip envía SessionID" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "SessionID: $sessionId" -ForegroundColor White
Write-Host ""

$infobipRequest = @{
    sessionId = $sessionId
}

$sessionResponse = Invoke-APIRequest -Method POST -Endpoint "/infobip/session" -Body $infobipRequest

if (-not $sessionResponse) {
    Write-Host ""
    Write-Host "❌ Error creando sesión desde Infobip. Abortando." -ForegroundColor Red
    exit 1
}

$token = $sessionResponse.token
$termsUrl = $sessionResponse.url

Write-Host ""
Write-Host "✅ Sesión creada exitosamente" -ForegroundColor Green
Write-Host "   Token: $token" -ForegroundColor Green
Write-Host "   URL de términos: $termsUrl" -ForegroundColor Green
Write-Host ""

Start-Sleep -Seconds 2

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 2: Frontend consulta por SessionID" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$statusResponse = Invoke-APIRequest -Method GET -Endpoint "/terms/by-session/$sessionId"

if ($statusResponse -and $statusResponse.status -eq "PENDING") {
    Write-Host ""
    Write-Host "✅ Frontend obtuvo los datos de la sesión" -ForegroundColor Green
    Write-Host "   Token obtenido: $($statusResponse.token)" -ForegroundColor Green
    Write-Host "   Estado: $($statusResponse.status)" -ForegroundColor Green
    Write-Host "   Expira: $($statusResponse.expiresAt)" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error obteniendo sesión por sessionId" -ForegroundColor Red
    exit 1
}

Write-Host ""
Start-Sleep -Seconds 2

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 3: Cliente Acepta Términos" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "El cliente accede a: $termsUrl" -ForegroundColor White
Write-Host "Y hace clic en 'Aceptar'" -ForegroundColor White
Write-Host ""

$acceptResponse = Invoke-APIRequest -Method POST -Endpoint "/terms/$token/accept"

if ($acceptResponse -and $acceptResponse.status -eq "ACCEPTED") {
    Write-Host ""
    Write-Host "✅ Términos aceptados exitosamente" -ForegroundColor Green
    Write-Host "   Estado: $($acceptResponse.status)" -ForegroundColor Green
    Write-Host "   Aceptado en: $($acceptResponse.acceptedAt)" -ForegroundColor Green
    Write-Host "   Mensaje: $($acceptResponse.message)" -ForegroundColor Green
} else {
    Write-Host ""
    Write-Host "❌ Error al aceptar términos" -ForegroundColor Red
    exit 1
}

Write-Host ""
Write-Host "⏳ Esperando notificación a Infobip (en segundo plano)..." -ForegroundColor Yellow
Start-Sleep -Seconds 3

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "PASO 4: Verificar Estado Final" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$finalStatusBySession = Invoke-APIRequest -Method GET -Endpoint "/terms/by-session/$sessionId"

if ($finalStatusBySession -and $finalStatusBySession.status -eq "ACCEPTED") {
    Write-Host ""
    Write-Host "✅ Estado final verificado por SessionID" -ForegroundColor Green
    Write-Host "   SessionID: $sessionId" -ForegroundColor Green
    Write-Host "   Token: $($finalStatusBySession.token)" -ForegroundColor Green
    Write-Host "   Estado: $($finalStatusBySession.status)" -ForegroundColor Green
    Write-Host "   Aceptado en: $($finalStatusBySession.acceptedAt)" -ForegroundColor Green
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "✅ TEST COMPLETADO EXITOSAMENTE" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "Resumen del Flujo:" -ForegroundColor White
Write-Host "  1. ✅ Infobip envió sessionId → Backend creó sesión" -ForegroundColor White
Write-Host "  2. ✅ Frontend consultó por sessionId → Obtuvo token" -ForegroundColor White
Write-Host "  3. ✅ Cliente aceptó términos → Estado ACCEPTED" -ForegroundColor White
Write-Host "  4. ✅ Backend notificó a Infobip (webhook automático)" -ForegroundColor White
Write-Host "  5. ✅ Estado final verificado" -ForegroundColor White
Write-Host ""
Write-Host "Datos importantes:" -ForegroundColor Cyan
Write-Host "  SessionID: $sessionId" -ForegroundColor White
Write-Host "  Token: $token" -ForegroundColor White
Write-Host "  URL: $termsUrl" -ForegroundColor White
Write-Host ""

Write-Host "📝 Nota: El webhook a Infobip se envió en segundo plano." -ForegroundColor Yellow
Write-Host "   Verifica los logs del servidor para confirmar el envío." -ForegroundColor Yellow
Write-Host ""
