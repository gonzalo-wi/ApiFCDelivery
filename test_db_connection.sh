#!/bin/bash
# Script para diagnosticar conectividad a la base de datos

echo "=== Test de Conectividad a Base de Datos ==="
echo ""

DB_HOST="192.168.0.227"
DB_PORT="3306"

echo "1. Verificando conectividad de red al host $DB_HOST..."
if ping -c 3 $DB_HOST; then
    echo "✓ Ping exitoso"
else
    echo "✗ No se puede hacer ping al host"
fi

echo ""
echo "2. Verificando puerto MySQL ($DB_PORT)..."
if command -v nc &> /dev/null; then
    if nc -zv $DB_HOST $DB_PORT 2>&1; then
        echo "✓ Puerto MySQL accesible"
    else
        echo "✗ Puerto MySQL no accesible"
    fi
elif command -v telnet &> /dev/null; then
    timeout 5 telnet $DB_HOST $DB_PORT
else
    echo "⚠ nc o telnet no disponibles para probar el puerto"
fi

echo ""
echo "3. Mostrando tabla de rutas..."
ip route | grep default

echo ""
echo "4. Mostrando interfaz de red..."
ip addr show | grep -E "inet |^[0-9]"

echo ""
echo "5. Verificando firewall (iptables)..."
if command -v iptables &> /dev/null; then
    sudo iptables -L -n | head -20
else
    echo "iptables no disponible"
fi

echo ""
echo "=== Fin del diagnóstico ==="
