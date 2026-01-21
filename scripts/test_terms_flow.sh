#!/bin/bash

# Script de prueba del flujo de Términos y Condiciones con Infobip
# Asegúrate de que el servidor esté corriendo en localhost:8080

echo "🧪 Iniciando pruebas del flujo de Términos y Condiciones"
echo "=========================================================="
echo ""

BASE_URL="http://localhost:8080/api/v1"

# Color codes
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Test 1: Crear sesión
echo -e "${YELLOW}Test 1: Crear sesión desde Infobip${NC}"
echo "POST $BASE_URL/infobip/session"

RESPONSE=$(curl -s -X POST "$BASE_URL/infobip/session" \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-session-'$(date +%s)'"}')

echo "$RESPONSE" | jq '.'

# Extraer token
TOKEN=$(echo "$RESPONSE" | jq -r '.token')

if [ "$TOKEN" != "null" ] && [ -n "$TOKEN" ]; then
    echo -e "${GREEN}✓ Token generado: ${TOKEN:0:20}...${NC}"
else
    echo -e "${RED}✗ Error: No se pudo generar token${NC}"
    exit 1
fi

echo ""
sleep 1

# Test 2: Consultar estado (debe estar PENDING)
echo -e "${YELLOW}Test 2: Consultar estado del token${NC}"
echo "GET $BASE_URL/terms/$TOKEN"

STATUS_RESPONSE=$(curl -s "$BASE_URL/terms/$TOKEN")
echo "$STATUS_RESPONSE" | jq '.'

STATUS=$(echo "$STATUS_RESPONSE" | jq -r '.status')

if [ "$STATUS" == "PENDING" ]; then
    echo -e "${GREEN}✓ Estado correcto: PENDING${NC}"
else
    echo -e "${RED}✗ Estado incorrecto: $STATUS (esperado: PENDING)${NC}"
fi

echo ""
sleep 1

# Test 3: Aceptar términos
echo -e "${YELLOW}Test 3: Aceptar términos${NC}"
echo "POST $BASE_URL/terms/$TOKEN/accept"

ACCEPT_RESPONSE=$(curl -s -X POST "$BASE_URL/terms/$TOKEN/accept" \
  -H "Content-Type: application/json" \
  -H "User-Agent: curl-test-script")

echo "$ACCEPT_RESPONSE" | jq '.'

ACCEPT_STATUS=$(echo "$ACCEPT_RESPONSE" | jq -r '.status')

if [ "$ACCEPT_STATUS" == "ACCEPTED" ]; then
    echo -e "${GREEN}✓ Términos aceptados correctamente${NC}"
else
    echo -e "${RED}✗ Error al aceptar términos${NC}"
fi

echo ""
sleep 1

# Test 4: Verificar idempotencia (aceptar nuevamente)
echo -e "${YELLOW}Test 4: Probar idempotencia (aceptar de nuevo)${NC}"
echo "POST $BASE_URL/terms/$TOKEN/accept"

IDEMPOTENT_RESPONSE=$(curl -s -X POST "$BASE_URL/terms/$TOKEN/accept" \
  -H "Content-Type: application/json")

echo "$IDEMPOTENT_RESPONSE" | jq '.'

MESSAGE=$(echo "$IDEMPOTENT_RESPONSE" | jq -r '.message')

if [[ "$MESSAGE" == *"previamente"* ]]; then
    echo -e "${GREEN}✓ Idempotencia funciona correctamente${NC}"
else
    echo -e "${RED}✗ Idempotencia no funcionó como esperado${NC}"
fi

echo ""
sleep 1

# Test 5: Consultar estado final (debe estar ACCEPTED)
echo -e "${YELLOW}Test 5: Consultar estado final${NC}"
echo "GET $BASE_URL/terms/$TOKEN"

FINAL_STATUS_RESPONSE=$(curl -s "$BASE_URL/terms/$TOKEN")
echo "$FINAL_STATUS_RESPONSE" | jq '.'

FINAL_STATUS=$(echo "$FINAL_STATUS_RESPONSE" | jq -r '.status')

if [ "$FINAL_STATUS" == "ACCEPTED" ]; then
    echo -e "${GREEN}✓ Estado final correcto: ACCEPTED${NC}"
else
    echo -e "${RED}✗ Estado final incorrecto: $FINAL_STATUS${NC}"
fi

echo ""
echo "=========================================================="
echo -e "${GREEN}🎉 Pruebas completadas${NC}"
echo ""

# Test 6: Crear una nueva sesión y rechazar
echo -e "${YELLOW}Test 6: Crear sesión y rechazar términos${NC}"

REJECT_RESPONSE=$(curl -s -X POST "$BASE_URL/infobip/session" \
  -H "Content-Type: application/json" \
  -d '{"sessionId": "test-reject-'$(date +%s)'"}')

REJECT_TOKEN=$(echo "$REJECT_RESPONSE" | jq -r '.token')

if [ "$REJECT_TOKEN" != "null" ] && [ -n "$REJECT_TOKEN" ]; then
    echo -e "${GREEN}✓ Token para rechazo generado${NC}"
    
    echo "POST $BASE_URL/terms/$REJECT_TOKEN/reject"
    
    REJECT_RESULT=$(curl -s -X POST "$BASE_URL/terms/$REJECT_TOKEN/reject" \
      -H "Content-Type: application/json")
    
    echo "$REJECT_RESULT" | jq '.'
    
    REJECT_STATUS=$(echo "$REJECT_RESULT" | jq -r '.status')
    
    if [ "$REJECT_STATUS" == "REJECTED" ]; then
        echo -e "${GREEN}✓ Términos rechazados correctamente${NC}"
    else
        echo -e "${RED}✗ Error al rechazar términos${NC}"
    fi
else
    echo -e "${RED}✗ No se pudo crear sesión para rechazo${NC}"
fi

echo ""
echo "=========================================================="
echo -e "${GREEN}✅ Todas las pruebas completadas${NC}"
echo ""
echo "Tokens generados para inspección manual:"
echo "  - Token aceptado: $TOKEN"
echo "  - Token rechazado: $REJECT_TOKEN"
