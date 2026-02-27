# Script de prueba para verificar que GET /deliveries devuelve entregado_por y tipos_dispensers
# Autor: GoFrioCalor API Testing
# Fecha: 2026-02-26

$BASE_URL = "http://localhost:8080/api/v1"

Write-Host "`n========================================" -ForegroundColor Cyan
Write-Host "TEST: GET Delivery Response Fields" -ForegroundColor Cyan
Write-Host "========================================`n" -ForegroundColor Cyan

# Paso 1: Crear una delivery desde Infobip con tipos P y M
Write-Host "[PASO 1] Creando delivery con tipos P:2 y M:1..." -ForegroundColor Yellow

$createBody = @{
    nro_cta = "TEST-$(Get-Random -Minimum 1000 -Maximum 9999)"
    nro_rto = "RTO-$(Get-Random -Minimum 1000 -Maximum 9999)"
    tipos = @{
        P = 2
        M = 1
    }
    tipo_entrega = "Instalacion"
    entregado_por = "Repartidor"
    session_id = "session-test-$(Get-Random)"
} | ConvertTo-Json

try {
    $createResponse = Invoke-RestMethod -Uri "$BASE_URL/deliveries/infobip" `
        -Method POST `
        -ContentType "application/json" `
        -Body $createBody

    Write-Host "✓ Delivery creada exitosamente" -ForegroundColor Green
    Write-Host "  Token: $($createResponse.token)" -ForegroundColor Gray
    
    $token = $createResponse.token
} catch {
    Write-Host "✗ Error al crear delivery: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Esperar un momento para que se guarde
Start-Sleep -Seconds 1

# Paso 2: Obtener la delivery por token (buscando en la lista)
Write-Host "`n[PASO 2] Obteniendo lista de deliveries..." -ForegroundColor Yellow

try {
    $deliveries = Invoke-RestMethod -Uri "$BASE_URL/deliveries" -Method GET
    
    # Buscar la delivery que acabamos de crear
    $delivery = $deliveries | Where-Object { $_.token -eq $token } | Select-Object -First 1
    
    if (-not $delivery) {
        Write-Host "✗ No se encontró la delivery con token $token" -ForegroundColor Red
        exit 1
    }
    
    Write-Host "✓ Delivery encontrada (ID: $($delivery.id))" -ForegroundColor Green
    
} catch {
    Write-Host "✗ Error al obtener deliveries: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
}

# Paso 3: Verificar campos en la respuesta
Write-Host "`n[PASO 3] Verificando campos en la respuesta..." -ForegroundColor Yellow
Write-Host "`n----------------------------------------" -ForegroundColor White

$allChecksPass = $true

# Verificar entregado_por
Write-Host "Campo 'entregado_por':" -NoNewline
if ($delivery.PSObject.Properties.Name -contains "entregado_por") {
    Write-Host " ✓ PRESENTE" -ForegroundColor Green
    Write-Host "  Valor: $($delivery.entregado_por)" -ForegroundColor Gray
} else {
    Write-Host " ✗ AUSENTE" -ForegroundColor Red
    $allChecksPass = $false
}

# Verificar tipos_dispensers
Write-Host "`nCampo 'tipos_dispensers':" -NoNewline
if ($delivery.PSObject.Properties.Name -contains "tipos_dispensers") {
    Write-Host " ✓ PRESENTE" -ForegroundColor Green
    
    $tipos = $delivery.tipos_dispensers
    
    # Verificar cada tipo
    Write-Host "  Tipos encontrados:" -ForegroundColor Gray
    if ($tipos.P) { Write-Host "    P: $($tipos.P)" -ForegroundColor Cyan }
    if ($tipos.M) { Write-Host "    M: $($tipos.M)" -ForegroundColor Cyan }
    if ($tipos.A) { Write-Host "    A: $($tipos.A)" -ForegroundColor Cyan }
    if ($tipos.B) { Write-Host "    B: $($tipos.B)" -ForegroundColor Cyan }
    if ($tipos.C) { Write-Host "    C: $($tipos.C)" -ForegroundColor Cyan }
    if ($tipos.HELADERA) { Write-Host "    HELADERA: $($tipos.HELADERA)" -ForegroundColor Cyan }
    
    # Verificar que los valores coincidan con lo enviado (P:2, M:1)
    if ($tipos.P -eq 2 -and $tipos.M -eq 1) {
        Write-Host "  ✓ Conteo correcto (P:2, M:1)" -ForegroundColor Green
    } else {
        Write-Host "  ✗ Conteo incorrecto. Esperado P:2, M:1. Recibido P:$($tipos.P), M:$($tipos.M)" -ForegroundColor Red
        $allChecksPass = $false
    }
} else {
    Write-Host " ✗ AUSENTE" -ForegroundColor Red
    $allChecksPass = $false
}

# Verificar dispensers array
Write-Host "`nCampo 'dispensers' (array):" -NoNewline
if ($delivery.PSObject.Properties.Name -contains "dispensers") {
    Write-Host " ✓ PRESENTE" -ForegroundColor Green
    Write-Host "  Cantidad de dispensers: $($delivery.dispensers.Count)" -ForegroundColor Gray
    
    if ($delivery.dispensers.Count -eq 3) {
        Write-Host "  ✓ Cantidad correcta (3 dispensers)" -ForegroundColor Green
        
        # Mostrar cada dispenser
        foreach ($disp in $delivery.dispensers) {
            Write-Host "    - Tipo: $($disp.tipo), Marca: $($disp.marca), Serie: $($disp.nro_serie)" -ForegroundColor Gray
        }
    } else {
        Write-Host "  ✗ Cantidad incorrecta. Esperado: 3, Recibido: $($delivery.dispensers.Count)" -ForegroundColor Red
        $allChecksPass = $false
    }
} else {
    Write-Host " ✗ AUSENTE" -ForegroundColor Red
    $allChecksPass = $false
}

# Paso 4: Mostrar JSON completo
Write-Host "`n[PASO 4] Respuesta JSON completa:" -ForegroundColor Yellow
Write-Host "----------------------------------------" -ForegroundColor White
$delivery | ConvertTo-Json -Depth 5
Write-Host "----------------------------------------" -ForegroundColor White

# Resultado final
Write-Host "`n========================================" -ForegroundColor Cyan
if ($allChecksPass) {
    Write-Host "✓ TODAS LAS VERIFICACIONES PASARON" -ForegroundColor Green
    Write-Host "========================================`n" -ForegroundColor Cyan
    exit 0
} else {
    Write-Host "✗ ALGUNAS VERIFICACIONES FALLARON" -ForegroundColor Red
    Write-Host "========================================`n" -ForegroundColor Cyan
    exit 1
}
