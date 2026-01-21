# Script de prueba del flujo de Términos y Condiciones con Infobip (PowerShell)
# Asegúrate de que el servidor esté corriendo en localhost:8080

Write-Host "🧪 Iniciando pruebas del flujo de Términos y Condiciones" -ForegroundColor Cyan
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host ""

$BASE_URL = "http://localhost:8080/api/v1"

# Test 1: Crear sesión
Write-Host "Test 1: Crear sesión desde Infobip" -ForegroundColor Yellow
Write-Host "POST $BASE_URL/infobip/session"

$sessionId = "test-session-$(Get-Date -Format 'yyyyMMddHHmmss')"
$body = @{
    sessionId = $sessionId
} | ConvertTo-Json

try {
    $response = Invoke-RestMethod -Uri "$BASE_URL/infobip/session" `
        -Method Post `
        -ContentType "application/json" `
        -Body $body

    $response | ConvertTo-Json -Depth 10 | Write-Host
    $TOKEN = $response.token

    if ($TOKEN) {
        Write-Host "✓ Token generado: $($TOKEN.Substring(0, 20))..." -ForegroundColor Green
    } else {
        Write-Host "✗ Error: No se pudo generar token" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "✗ Error en la petición: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

Write-Host ""
Start-Sleep -Seconds 1

# Test 2: Consultar estado (debe estar PENDING)
Write-Host "Test 2: Consultar estado del token" -ForegroundColor Yellow
Write-Host "GET $BASE_URL/terms/$TOKEN"

try {
    $statusResponse = Invoke-RestMethod -Uri "$BASE_URL/terms/$TOKEN" -Method Get
    $statusResponse | ConvertTo-Json -Depth 10 | Write-Host

    if ($statusResponse.status -eq "PENDING") {
        Write-Host "✓ Estado correcto: PENDING" -ForegroundColor Green
    } else {
        Write-Host "✗ Estado incorrecto: $($statusResponse.status) (esperado: PENDING)" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error consultando estado: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Start-Sleep -Seconds 1

# Test 3: Aceptar términos
Write-Host "Test 3: Aceptar términos" -ForegroundColor Yellow
Write-Host "POST $BASE_URL/terms/$TOKEN/accept"

try {
    $acceptResponse = Invoke-RestMethod -Uri "$BASE_URL/terms/$TOKEN/accept" `
        -Method Post `
        -ContentType "application/json" `
        -Headers @{"User-Agent" = "PowerShell-test-script"}

    $acceptResponse | ConvertTo-Json -Depth 10 | Write-Host

    if ($acceptResponse.status -eq "ACCEPTED") {
        Write-Host "✓ Términos aceptados correctamente" -ForegroundColor Green
    } else {
        Write-Host "✗ Error al aceptar términos" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error aceptando términos: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Start-Sleep -Seconds 1

# Test 4: Verificar idempotencia (aceptar nuevamente)
Write-Host "Test 4: Probar idempotencia (aceptar de nuevo)" -ForegroundColor Yellow
Write-Host "POST $BASE_URL/terms/$TOKEN/accept"

try {
    $idempotentResponse = Invoke-RestMethod -Uri "$BASE_URL/terms/$TOKEN/accept" `
        -Method Post `
        -ContentType "application/json"

    $idempotentResponse | ConvertTo-Json -Depth 10 | Write-Host

    if ($idempotentResponse.message -like "*previamente*") {
        Write-Host "✓ Idempotencia funciona correctamente" -ForegroundColor Green
    } else {
        Write-Host "✗ Idempotencia no funcionó como esperado" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error en prueba de idempotencia: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Start-Sleep -Seconds 1

# Test 5: Consultar estado final (debe estar ACCEPTED)
Write-Host "Test 5: Consultar estado final" -ForegroundColor Yellow
Write-Host "GET $BASE_URL/terms/$TOKEN"

try {
    $finalStatusResponse = Invoke-RestMethod -Uri "$BASE_URL/terms/$TOKEN" -Method Get
    $finalStatusResponse | ConvertTo-Json -Depth 10 | Write-Host

    if ($finalStatusResponse.status -eq "ACCEPTED") {
        Write-Host "✓ Estado final correcto: ACCEPTED" -ForegroundColor Green
    } else {
        Write-Host "✗ Estado final incorrecto: $($finalStatusResponse.status)" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error consultando estado final: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "🎉 Pruebas de aceptación completadas" -ForegroundColor Green
Write-Host ""

# Test 6: Crear una nueva sesión y rechazar
Write-Host "Test 6: Crear sesión y rechazar términos" -ForegroundColor Yellow

$rejectSessionId = "test-reject-$(Get-Date -Format 'yyyyMMddHHmmss')"
$rejectBody = @{
    sessionId = $rejectSessionId
} | ConvertTo-Json

try {
    $rejectResponse = Invoke-RestMethod -Uri "$BASE_URL/infobip/session" `
        -Method Post `
        -ContentType "application/json" `
        -Body $rejectBody

    $REJECT_TOKEN = $rejectResponse.token

    if ($REJECT_TOKEN) {
        Write-Host "✓ Token para rechazo generado" -ForegroundColor Green
        Write-Host "POST $BASE_URL/terms/$REJECT_TOKEN/reject"

        $rejectResult = Invoke-RestMethod -Uri "$BASE_URL/terms/$REJECT_TOKEN/reject" `
            -Method Post `
            -ContentType "application/json"

        $rejectResult | ConvertTo-Json -Depth 10 | Write-Host

        if ($rejectResult.status -eq "REJECTED") {
            Write-Host "✓ Términos rechazados correctamente" -ForegroundColor Green
        } else {
            Write-Host "✗ Error al rechazar términos" -ForegroundColor Red
        }
    } else {
        Write-Host "✗ No se pudo crear sesión para rechazo" -ForegroundColor Red
    }
} catch {
    Write-Host "✗ Error en prueba de rechazo: $($_.Exception.Message)" -ForegroundColor Red
}

Write-Host ""
Write-Host "==========================================================" -ForegroundColor Cyan
Write-Host "✅ Todas las pruebas completadas" -ForegroundColor Green
Write-Host ""
Write-Host "Tokens generados para inspección manual:"
Write-Host "  - Token aceptado: $TOKEN"
Write-Host "  - Token rechazado: $REJECT_TOKEN"
