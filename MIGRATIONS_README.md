# Sistema de Migraciones Automáticas

## 🚀 Funcionamiento

El proyecto ahora incluye un sistema de migraciones automático que se ejecuta cada vez que levantes el contenedor o inicies la aplicación. 

### ¿Cómo funciona?

1. **Al iniciar la aplicación**: Se ejecutan automáticamente todas las migraciones SQL pendientes
2. **Tracking**: Se crea una tabla `migrations` que registra qué migraciones ya se aplicaron
3. **Idempotente**: Puedes ejecutar la aplicación múltiples veces sin problemas - solo aplica migraciones nuevas

## 📁 Estructura

```
migrations/
├── 001_create_terms_sessions.sql
├── 002_add_session_id_to_deliveries.sql
├── 003_add_dispenser_types_p_m.sql
├── 003b_cleanup_duplicate_sessions.sql
├── 004_add_unique_session_id.sql
└── 005_add_client_info_to_deliveries.sql
```

Las migraciones se ejecutan **en orden alfabético** por nombre de archivo.

## ✅ Flujo de Inicialización

```
1. Aplicación inicia
2. Conexión a base de datos ✓
3. GORM AutoMigrate (crea tablas básicas) ✓
4. MigrationService ejecuta archivos SQL en orden:
   - Verifica si la migración ya fue aplicada
   - Si no, la ejecuta y la marca como aplicada
   - Continúa con la siguiente
5. Aplicación lista para usar ✓
```

## 🆕 Agregar Nueva Migración

Para agregar una nueva migración, simplemente:

1. Crea un archivo en `migrations/` con nombre ordenado:
   ```
   006_descripcion_del_cambio.sql
   ```

2. Escribe tu SQL:
   ```sql
   -- Descripción de la migración
   ALTER TABLE deliveries 
   ADD COLUMN nuevo_campo VARCHAR(100);
   
   CREATE INDEX idx_nuevo_campo ON deliveries(nuevo_campo);
   ```

3. **¡Listo!** La próxima vez que levantes la aplicación, se aplicará automáticamente.

## 🔒 Seguridad

- Cada migración se ejecuta en una **transacción**
- Si una migración falla, se hace rollback automático
- Los errores de "columna ya existe" se ignoran automáticamente (idempotencia)

## 🐳 Con Docker

Cuando ejecutes:

```bash
docker-compose up --build
```

El sistema:
1. Levanta el contenedor
2. Conecta a la base de datos
3. Ejecuta TODAS las migraciones automáticamente
4. Inicia el servidor

**No necesitas ejecutar scripts manualmente.**

## 📊 Tabla de Tracking

Se crea automáticamente una tabla `migrations`:

| Campo      | Tipo    | Descripción                          |
|------------|---------|--------------------------------------|
| id         | UINT    | ID autoincremental                   |
| name       | STRING  | Nombre del archivo .sql              |
| applied_at | INT64   | Timestamp de cuándo se aplicó        |

Ejemplo de contenido:
```
id | name                                | applied_at   
1  | 001_create_terms_sessions.sql       | 1709136000
2  | 002_add_session_id_to_deliveries.sql| 1709136001
3  | 003_add_dispenser_types_p_m.sql     | 1709136002
```

## 🧪 Base de Datos Nueva

Si creas una base de datos completamente nueva:

1. GORM crea las tablas básicas: `deliveries`, `dispensers`, `work_orders`, `terms_sessions`
2. El MigrationService aplica TODAS las migraciones SQL en orden
3. Resultado: Base de datos completamente configurada y lista

## 🛠️ Troubleshooting

### Ver qué migraciones se aplicaron

```sql
SELECT * FROM migrations ORDER BY applied_at;
```

### Re-ejecutar una migración

Si necesitas re-ejecutar una migración (por ejemplo, la modificaste):

```sql
-- Eliminar el registro
DELETE FROM migrations WHERE name = '005_add_client_info_to_deliveries.sql';
```

Luego reinicia la aplicación y se volverá a ejecutar.

### Error: "Migration already applied"

Esto es normal - significa que la migración ya se ejecutó antes. No pasa nada.

## 📝 Logs

Cuando inicies la aplicación, verás logs como:

```
INF Database connected successfully
INF Found 6 migration files
INF Applying migration: 001_create_terms_sessions.sql
INF Migration 001_create_terms_sessions.sql applied successfully
INF Migration 002_add_session_id_to_deliveries.sql already applied, skipping
...
INF All migrations completed successfully
INF Database migrations completed successfully
```

## 🎯 Beneficios

✅ **Automatización completa**: No más scripts manuales  
✅ **Reproducible**: Cualquier entorno nuevo se configura igual  
✅ **Versionado**: Las migraciones están en Git  
✅ **Seguro**: Transacciones y manejo de errores  
✅ **Idempotente**: Ejecuta múltiples veces sin problemas  

## 🔄 Workflow de Desarrollo

```bash
# 1. Crear nueva migración
echo "ALTER TABLE deliveries ADD COLUMN test VARCHAR(50);" > migrations/006_test.sql

# 2. Rebuild y levantar
docker-compose down
docker-compose up --build

# ✅ La migración se aplica automáticamente
```

## 🚫 Scripts Obsoletos

Ya no necesitas ejecutar manualmente:
- ❌ `apply_migration_002.ps1`
- ❌ `apply_migration_004.ps1`
- ❌ `apply_migration_005.ps1`

Todo se ejecuta automáticamente en el orden correcto.
