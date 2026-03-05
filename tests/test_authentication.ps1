# Script de prueba para autenticación JWT
# Asegúrate de que el servidor esté corriendo antes de ejecutar este script

Write-Host "=== Prueba de Autenticación JWT ===" -ForegroundColor Cyan
Write-Host ""

$apiBaseURL = "http://localhost:8095"

# 1. Obtener Token
Write-Host "1. Obteniendo token JWT..." -ForegroundColor Yellow
$headers = @{
    "x-api-key" = "MOBEUS_kG7pX2sV9nQ4aJ1cL8rT0yZ5wH3eU6mF2dC9bA1sR4xP"
}

try {
    $tokenResponse = Invoke-RestMethod -Uri "$apiBaseURL/dispenser-operations/auth/generar-token" -Method GET -Headers $headers
    $token = $tokenResponse.token
    Write-Host "✓ Token obtenido exitosamente" -ForegroundColor Green
    Write-Host "  Token: $($token.Substring(0, 50))..." -ForegroundColor Gray
    Write-Host "  Proveedor: $($tokenResponse.proveedor)" -ForegroundColor Gray
    Write-Host "  Expira: $($tokenResponse.expires_at_local)" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "✗ Error obteniendo token: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# 2. Probar sin token (debe fallar)
Write-Host "2. Probando endpoint sin token (debe fallar con 401)..." -ForegroundColor Yellow
try {
    $response = Invoke-RestMethod -Uri "$apiBaseURL/dispenser-operations/api/v1/work-orders" -Method GET
    Write-Host "✗ ERROR: La request debería haber fallado sin token" -ForegroundColor Red
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        Write-Host "✓ Correcto: Request rechazada con 401 Unauthorized" -ForegroundColor Green
    } else {
        Write-Host "✗ Error inesperado: Status code $statusCode" -ForegroundColor Red
    }
}
Write-Host ""

# 3. Probar con token inválido
Write-Host "3. Probando con token inválido (debe fallar con 401)..." -ForegroundColor Yellow
$invalidHeaders = @{
    "Authorization" = "Bearer token_invalido_12345"
    "Content-Type" = "application/json"
}

try {
    $response = Invoke-RestMethod -Uri "$apiBaseURL/dispenser-operations/api/v1/work-orders" -Method GET -Headers $invalidHeaders
    Write-Host "✗ ERROR: La request debería haber fallado con token inválido" -ForegroundColor Red
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        Write-Host "✓ Correcto: Token inválido rechazado con 401 Unauthorized" -ForegroundColor Green
    } else {
        Write-Host "✗ Error inesperado: Status code $statusCode" -ForegroundColor Red
    }
}
Write-Host ""

# 4. Probar con token válido
Write-Host "4. Probando con token válido..." -ForegroundColor Yellow
$validHeaders = @{
    "Authorization" = "Bearer $token"
    "Content-Type" = "application/json"
}

try {
    $response = Invoke-RestMethod -Uri "$apiBaseURL/dispenser-operations/api/v1/work-orders" -Method GET -Headers $validHeaders
    Write-Host "✓ Request exitosa con token válido" -ForegroundColor Green
    Write-Host "  Respuesta: $($response | ConvertTo-Json -Compress -Depth 2)" -ForegroundColor Gray
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "✗ Error: Status code $statusCode" -ForegroundColor Red
    Write-Host "  Mensaje: $($_.Exception.Message)" -ForegroundColor Gray
}
Write-Host ""

# 5. Verificar endpoint público (health check)
Write-Host "5. Verificando endpoint público /health (sin token)..." -ForegroundColor Yellow
try {
    $healthResponse = Invoke-RestMethod -Uri "$apiBaseURL/health" -Method GET
    Write-Host "✓ Health check accesible sin autenticación" -ForegroundColor Green
    Write-Host "  Status: $($healthResponse.status)" -ForegroundColor Gray
} catch {
    Write-Host "✗ Error accediendo al health check: $($_.Exception.Message)" -ForegroundColor Red
}
Write-Host ""

# 6. Probar formato incorrecto de Authorization
Write-Host "6. Probando con formato incorrecto de Authorization..." -ForegroundColor Yellow
$incorrectHeaders = @{
    "Authorization" = $token  # Sin "Bearer "
    "Content-Type" = "application/json"
}

try {
    $response = Invoke-RestMethod -Uri "$apiBaseURL/dispenser-operations/api/v1/work-orders" -Method GET -Headers $incorrectHeaders
    Write-Host "✗ ERROR: Debería haber rechazado el formato incorrecto" -ForegroundColor Red
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    if ($statusCode -eq 401) {
        Write-Host "✓ Correcto: Formato inválido rechazado con 401 Unauthorized" -ForegroundColor Green
    } else {
        Write-Host "✗ Error inesperado: Status code $statusCode" -ForegroundColor Red
    }
}
Write-Host ""

Write-Host "=== Pruebas Completadas ===" -ForegroundColor Cyan
Write-Host ""
Write-Host "Resumen:" -ForegroundColor Yellow
Write-Host "- Los endpoints de la API están protegidos con JWT" -ForegroundColor White
Write-Host "- Los tokens se obtienen del servicio de autenticación externo" -ForegroundColor White
Write-Host "- El endpoint /health permanece público" -ForegroundColor White
Write-Host "- Los tokens inválidos o mal formateados son rechazados" -ForegroundColor White
