# 🚀 Guía de Despliegue - Cambios de Estructura de Base de Datos

## 📋 Cambios en Este Release

### Nuevas Características
1. **Endpoint Infobip Delivery** (`POST /api/v1/deliveries/infobip`)
2. **Nuevos tipos de dispensers:**
   - `P` = Pie (Dispenser de pie)
   - `M` = Mesada (Dispenser de mesada)
3. **Helpers modularizados** para mejor mantenibilidad

### Cambios en Base de Datos
- Campo `tipo` en tabla `dispensers` ahora acepta valores adicionales: `'P'` y `'M'`
- **GORM AutoMigrate** maneja esto automáticamente ✅

---

## 🔄 Proceso de Despliegue

### Opción 1: Despliegue Automático con Docker (Recomendado)

#### 1. Hacer push de los cambios
```bash
git add .
git commit -m "feat: endpoint Infobip con tipos de dispensers P y M"
git push origin main
```

#### 2. En el servidor, actualizar el código
```bash
cd /ruta/a/tu/proyecto
git pull origin main
```

#### 3. Reconstruir y reiniciar el contenedor
```bash
# Detener contenedor actual
docker-compose down

# Reconstruir imagen con los nuevos cambios
docker-compose build --no-cache

# Iniciar contenedor
docker-compose up -d

# Ver logs para confirmar que arrancó bien
docker-compose logs -f app
```

#### 4. Verificar que GORM aplicó los cambios
```bash
# Ver logs del arranque
docker-compose logs app | grep "AutoMigrate"
```

Deberías ver algo como:
```
Conexión a la base de datos exitosa
Base de datos conectada y tablas migradas correctamente
```

---

### Opción 2: Despliegue Sin Interrupciones (Blue-Green)

```bash
# 1. Pull del código
git pull origin main

# 2. Construir nueva imagen con tag
docker build -t gofrocalor-api:new .

# 3. Detener contenedor viejo
docker-compose down

# 4. Actualizar docker-compose.yml para usar la nueva imagen
# (o simplemente hacer docker-compose up con rebuild)

# 5. Iniciar nuevo contenedor
docker-compose up -d

# 6. Verificar health check
curl http://localhost:8095/health
```

---

## 🗄️ Validación de Cambios en Base de Datos

### Verificar que los nuevos tipos funcionan

#### 1. Conectarse a la base de datos
```bash
# Desde el servidor
docker exec -it nombre-contenedor-mysql mysql -u usuario -p nombre_bd
```

#### 2. Verificar estructura de tabla dispensers
```sql
-- Ver el tipo de dato de la columna
DESCRIBE dispensers;

-- Ver tipos existentes
SELECT DISTINCT tipo FROM dispensers;
```

#### 3. Probar inserción manual (opcional)
```sql
-- Insertar un dispenser de tipo P para probar
INSERT INTO dispensers (marca, nro_serie, tipo, delivery_id, created_at, updated_at)
VALUES ('PENDIENTE', 'P-TEST-1', 'P', 1, NOW(), NOW());

-- Verificar
SELECT * FROM dispensers WHERE tipo = 'P';

-- Limpiar test
DELETE FROM dispensers WHERE nro_serie = 'P-TEST-1';
```

---

## 🧪 Testing Post-Despliegue

### 1. Health Check
```bash
curl http://tu-servidor:8095/health
```

**Respuesta esperada:**
```json
{
  "status": "ok",
  "timestamp": 1709000000
}
```

### 2. Probar el endpoint nuevo
```bash
curl -X POST http://tu-servidor:8095/api/v1/deliveries/infobip \
  -H "Content-Type: application/json" \
  -d '{
    "nro_cta": "CTA-TEST",
    "nro_rto": "RTO-TEST",
    "tipos": {
      "P": 2,
      "M": 1
    },
    "tipo_entrega": "Instalacion",
    "entregado_por": "Repartidor",
    "session_id": "TEST-DEPLOY-001"
  }'
```

**Respuesta esperada (201):**
```json
{
  "token": "1234",
  "message": "Entrega creada exitosamente"
}
```

### 3. Verificar en la base de datos
```sql
-- Ver la última entrega creada
SELECT * FROM deliveries ORDER BY id DESC LIMIT 1;

-- Ver los dispensers asociados
SELECT * FROM dispensers WHERE delivery_id = (
    SELECT id FROM deliveries ORDER BY id DESC LIMIT 1
);
```

Deberías ver dispensers con tipo 'P' y 'M'.

---

## 🔧 Troubleshooting

### Problema: El contenedor no arranca

**Solución:**
```bash
# Ver logs completos
docker-compose logs app

# Verificar errores de compilación
docker-compose build

# Verificar conectividad a BD
docker-compose exec app ping mysql-host
```

### Problema: Error "invalid tipo value"

**Causa:** La base de datos tiene una restricción ENUM antigua.

**Solución:**
```sql
-- Conectarse a MySQL
mysql -u usuario -p

-- Cambiar a VARCHAR si tiene ENUM
ALTER TABLE dispensers MODIFY COLUMN tipo VARCHAR(20);
```

### Problema: GORM no migra automáticamente

**Solución 1 - Forzar recreación:**
```bash
# CUIDADO: Esto borra los datos
docker-compose down -v  # Borra volúmenes
docker-compose up -d
```

**Solución 2 - Migración manual:**
```sql
-- Verificar que la columna acepta los nuevos valores
ALTER TABLE dispensers MODIFY COLUMN tipo VARCHAR(20) NOT NULL;
```

---

## 📊 Monitoreo Post-Despliegue

### Verificar logs en tiempo real
```bash
docker-compose logs -f app
```

### Métricas a observar
```bash
# Uso de CPU/RAM del contenedor
docker stats gofrocalor-api

# Conexiones a BD
docker-compose exec app netstat -an | grep 3306

# Requests por minuto (si tienes metrics)
# Ver logs de Gin/GIN-debug
```

### Verificar endpoints existentes no afectados
```bash
# GET deliveries
curl http://tu-servidor:8095/api/v1/deliveries

# POST delivery tradicional
curl -X POST http://tu-servidor:8095/api/v1/deliveries \
  -H "Content-Type: application/json" \
  -d '{...}'
```

---

## 🔐 Checklist de Seguridad

- [ ] Variables de entorno actualizadas en el servidor
- [ ] Firewall permite puerto 8095
- [ ] Backup de base de datos realizado antes del deploy
- [ ] CORS configurado correctamente
- [ ] Logs de errores monitoreados

---

## 📝 Notas Importantes

### ⚠️ Sobre GORM AutoMigrate

GORM AutoMigrate:
- ✅ **Agrega** columnas nuevas automáticamente
- ✅ **Modifica** tipos de datos si es seguro
- ✅ **Crea** índices nuevos
- ❌ **NO** elimina columnas (por seguridad)
- ❌ **NO** modifica constraints complejos

En nuestro caso:
- ✅ El cambio de validación `oneof=A B C HELADERA` a `oneof=A B C HELADERA P M` 
  es manejado a nivel de aplicación por GORM
- ✅ No requiere ALTER TABLE explícito
- ✅ La columna ya es VARCHAR, acepta cualquier valor

### 🔄 Rollback Plan

Si algo sale mal:

```bash
# 1. Revertir al commit anterior
git revert HEAD

# 2. Rebuild del contenedor
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# 3. O restaurar desde backup de BD si es necesario
mysql -u usuario -p nombre_bd < backup_pre_deploy.sql
```

---

## ✅ Validación Final

Después del despliegue, confirma:

1. ✅ Contenedor corriendo: `docker ps | grep gofrocalor`
2. ✅ Health check OK: `curl http://servidor:8095/health`
3. ✅ Endpoint Infobip funcional: probar con Postman
4. ✅ Endpoints anteriores funcionan: probar GET/POST deliveries
5. ✅ Logs sin errores: `docker-compose logs app | grep ERROR`
6. ✅ Dispensers con tipo P y M se crean correctamente

---

## 📞 Soporte

Si encuentras problemas:
1. Revisa logs: `docker-compose logs app`
2. Verifica conectividad BD
3. Consulta `docs/REFACTORING_SUMMARY.md` para detalles técnicos
4. Revisa `postman/` para ejemplos de testing
