# Script de prueba para el flujo Mobile + RabbitMQ
# Asegurate de tener el servidor corriendo antes de ejecutar

$BASE_URL = "http://localhost:8080/api/v1"
$headers = @{
    "Content-Type" = "application/json"
}

Write-Host "Iniciando prueba del flujo Mobile + RabbitMQ" -ForegroundColor Cyan
Write-Host ""

# Paso 1: Crear una entrega de prueba (simular Infobip)
Write-Host "Paso 1: Creando delivery de prueba..." -ForegroundColor Yellow
$deliveryData = @{
    nro_cta = "TEST12345"
    name = "Gonzalo Wiñazki"
    address = "Santiago de Liniers 3118"
    locality = "Ciudadela"
    nro_rto = "RTO999"
    cantidad = 2
    estado = "Pendiente"
    tipo_entrega = "Instalacion"
    entregado_por = "Tecnico"
    fecha_accion = "2025-11-12"
} | ConvertTo-Json

try {
    $delivery = Invoke-RestMethod -Uri "$BASE_URL/deliveries" -Method Post -Body $deliveryData -Headers $headers
    $deliveryId = $delivery.id
    $token = $delivery.token
    Write-Host "OK Delivery creado - ID: $deliveryId, Token: $token" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "ERROR Error creando delivery: $_" -ForegroundColor Red
    exit 1
}

# Paso 2: Agregar dispensers al delivery
Write-Host "Paso 2: Agregando dispensers..." -ForegroundColor Yellow
$dispenser1 = @{
    marca = "LAMO"
    nro_serie = "TEST-LM-001"
    tipo = "P"
    delivery_id = $deliveryId
} | ConvertTo-Json

$dispenser2 = @{
    marca = "LAMO"
    nro_serie = "TEST-LM-002"
    tipo = "M"
    delivery_id = $deliveryId
} | ConvertTo-Json

try {
    Invoke-RestMethod -Uri "$BASE_URL/dispensers" -Method Post -Body $dispenser1 -Headers $headers | Out-Null
    Invoke-RestMethod -Uri "$BASE_URL/dispensers" -Method Post -Body $dispenser2 -Headers $headers | Out-Null
    Write-Host "OK Dispensers agregados" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "ERROR Error agregando dispensers: $_" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# Paso 3: Simular App Movil - Validar Token
Write-Host "Paso 3: App Movil - Validando token del cliente..." -ForegroundColor Yellow
$validateTokenData = @{
    token = $token
    nro_cta = "TEST12345"
    fecha_accion = "2025-11-12"
} | ConvertTo-Json

try {
    $tokenResponse = Invoke-RestMethod -Uri "$BASE_URL/mobile/validate-token" -Method Post -Body $validateTokenData -Headers $headers
    
    if ($tokenResponse.valid) {
        Write-Host "OK Token valido!" -ForegroundColor Green
        Write-Host "   - Delivery ID: $($tokenResponse.delivery.id)" -ForegroundColor Gray
        Write-Host "   - Nro Cuenta: $($tokenResponse.delivery.nro_cta)" -ForegroundColor Gray
        Write-Host "   - Tipo: $($tokenResponse.delivery.tipo_entrega)" -ForegroundColor Gray
        Write-Host "   - Dispensers a escanear: $($tokenResponse.dispensers.Count)" -ForegroundColor Gray
        Write-Host ""
    } else {
        Write-Host "ERROR Token invalido" -ForegroundColor Red
        exit 1
    }
} catch {
    Write-Host "ERROR Error validando token: $_" -ForegroundColor Red
    exit 1
}

Start-Sleep -Seconds 1

# Paso 4: Simular escaneo de dispensers
Write-Host "Paso 4: Escaneando dispensers..." -ForegroundColor Yellow
$dispensers = @("TEST-LM-001", "TEST-LM-002")
$validatedDispensers = @()

foreach ($nroSerie in $dispensers) {
    $validateDispenserData = @{
        delivery_id = $deliveryId
        nro_serie = $nroSerie
    } | ConvertTo-Json
    
    try {
        $dispenserResponse = Invoke-RestMethod -Uri "$BASE_URL/mobile/validate-dispenser" -Method Post -Body $validateDispenserData -Headers $headers
        
        if ($dispenserResponse.valid) {
            Write-Host "   OK Dispenser $nroSerie validado" -ForegroundColor Green
            $validatedDispensers += $nroSerie
        } else {
            Write-Host "   ERROR Dispenser $nroSerie NO es valido" -ForegroundColor Red
        }
    } catch {
        Write-Host "   ERROR Error validando dispenser $nroSerie : $_" -ForegroundColor Red
    }
}

Write-Host ""
Start-Sleep -Seconds 1

# Paso 5: Completar la entrega
Write-Host "Paso 5: Completando la entrega..." -ForegroundColor Yellow
$completeData = @{
    delivery_id = $deliveryId
    token = $token
    validated_dispensers = $validatedDispensers
} | ConvertTo-Json

try {
    $completeResponse = Invoke-RestMethod -Uri "$BASE_URL/mobile/complete-delivery" -Method Post -Body $completeData -Headers $headers
    
    Write-Host "OK Entrega completada exitosamente!" -ForegroundColor Green
    Write-Host "   - Nro Cta: $($completeResponse.nroCta)" -ForegroundColor Gray
    Write-Host "   - Nombre: $($completeResponse.name)" -ForegroundColor Gray
    Write-Host "   - Direccion: $($completeResponse.address)" -ForegroundColor Gray
    Write-Host "   - Localidad: $($completeResponse.locality)" -ForegroundColor Gray
    Write-Host "   - Token: $($completeResponse.token)" -ForegroundColor Gray
    Write-Host "   - Dispensers: $($completeResponse.dispensers.Count)" -ForegroundColor Gray
    
    Write-Host ""
    Write-Host "Mensaje publicado a RabbitMQ en cola 'q.workorder.generate'" -ForegroundColor Cyan
    Write-Host "   El worker procesara el mensaje y:" -ForegroundColor Gray
    Write-Host "   1. Creara la orden de trabajo (WorkOrder)" -ForegroundColor Gray
    Write-Host "   2. Generara el PDF" -ForegroundColor Gray
    Write-Host "   3. Enviara el email al cliente" -ForegroundColor Gray
    Write-Host "   4. (Futuro) Guardara en storage" -ForegroundColor Gray
    Write-Host ""
} catch {
    Write-Host "ERROR Error completando entrega: $_" -ForegroundColor Red
    $errorDetails = $_.ErrorDetails.Message | ConvertFrom-Json
    if ($errorDetails.error) {
        Write-Host "   Detalle: $($errorDetails.error)" -ForegroundColor Red
    }
    exit 1
}

Write-Host ""
Write-Host "Prueba completada exitosamente!" -ForegroundColor Cyan
Write-Host "Verifica los logs del servidor para ver el procesamiento del worker" -ForegroundColor Yellow
Write-Host ""

# Opcional: Verificar en RabbitMQ
Write-Host "Puedes verificar RabbitMQ en: http://192.168.0.250:15672" -ForegroundColor Magenta
Write-Host "   Usuario: admin-  |  Password: admin123" -ForegroundColor Magenta
Write-Host ""
