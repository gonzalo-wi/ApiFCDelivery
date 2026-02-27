-- Script de preparación: Identificar y limpiar session_ids duplicados antes de aplicar índice único
-- Fecha: 2026-02-27
-- IMPORTANTE: Ejecutar ANTES de la migración 004

-- 1. Ver qué session_ids están duplicados
SELECT session_id, COUNT(*) as count, GROUP_CONCAT(id ORDER BY id) as delivery_ids
FROM deliveries 
WHERE session_id IS NOT NULL AND session_id != ''
GROUP BY session_id 
HAVING COUNT(*) > 1;

-- 2. Para cada duplicado, conservar solo el más antiguo (menor ID) y eliminar los demás
-- Opción A: Ver los IDs que serían eliminados (SIN ELIMINAR todavía)
SELECT d2.id, d2.session_id, d2.created_at
FROM deliveries d2
INNER JOIN (
    SELECT session_id, MIN(id) as min_id
    FROM deliveries
    WHERE session_id IS NOT NULL AND session_id != ''
    GROUP BY session_id
    HAVING COUNT(*) > 1
) d1 ON d2.session_id = d1.session_id AND d2.id > d1.min_id
ORDER BY d2.session_id, d2.id;

-- Opción B: ELIMINAR los duplicados (conservando el más antiguo)
-- DESCOMENTAR LA SIGUIENTE LÍNEA SOLO SI ESTÁS SEGURO:
-- DELETE d2
-- FROM deliveries d2
-- INNER JOIN (
--     SELECT session_id, MIN(id) as min_id
--     FROM deliveries
--     WHERE session_id IS NOT NULL AND session_id != ''
--     GROUP BY session_id
--     HAVING COUNT(*) > 1
-- ) d1 ON d2.session_id = d1.session_id AND d2.id > d1.min_id;

-- 3. Opción C: Si prefieres mantener todos los registros, puedes limpiar los session_id duplicados
-- UPDATE deliveries SET session_id = NULL WHERE id IN (...ids de duplicados...);
