# Script de prueba para el endpoint del contact center
# Este endpoint NO requiere autenticación

$baseUrl = "http://localhost:9090/dispenser-operations/api/v1"

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "Test: Contact Center - Obtener Token" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

# Parámetros de búsqueda
$fechaAccion = "2026-03-21"
$nroCta = "43534"

Write-Host "Buscando delivery con:" -ForegroundColor Yellow
Write-Host "  - Fecha Acción: $fechaAccion" -ForegroundColor White
Write-Host "  - Nro Cuenta: $nroCta" -ForegroundColor White
Write-Host ""

# Construir URL con query parameters
$url = "$baseUrl/deliveries/contact-center/token?fecha_accion=$fechaAccion&nro_cta=$nroCta"

try {
    Write-Host "Realizando petición GET a:" -ForegroundColor Yellow
    Write-Host "  $url" -ForegroundColor Gray
    Write-Host ""
    
    $response = Invoke-RestMethod -Uri $url -Method Get -ContentType "application/json"
    
    Write-Host "✓ Delivery encontrado!" -ForegroundColor Green
    Write-Host ""
    Write-Host "Respuesta:" -ForegroundColor Cyan
    $response | ConvertTo-Json -Depth 3 | Write-Host -ForegroundColor White
    Write-Host ""
    Write-Host "========================================" -ForegroundColor Green
    Write-Host "Token: $($response.token)" -ForegroundColor Yellow
    Write-Host "========================================" -ForegroundColor Green
    
} catch {
    $statusCode = $_.Exception.Response.StatusCode.value__
    Write-Host "✗ Error en la petición" -ForegroundColor Red
    Write-Host "Status Code: $statusCode" -ForegroundColor Red
    
    if ($_.ErrorDetails.Message) {
        Write-Host "Detalles del error:" -ForegroundColor Yellow
        $_.ErrorDetails.Message | ConvertFrom-Json | ConvertTo-Json -Depth 3 | Write-Host -ForegroundColor White
    }
}

Write-Host ""
Write-Host "Notas:" -ForegroundColor Cyan
Write-Host "  - Este endpoint NO requiere autenticación (x-api-key)" -ForegroundColor White
Write-Host "  - Formato de fecha: YYYY-MM-DD" -ForegroundColor White
Write-Host "  - Ambos parámetros son obligatorios" -ForegroundColor White
Write-Host ""
