-- Migration 009: Add validated_dispensers field to deliveries table
-- Description: Adds validated_dispensers JSONB column to store actual dispenser codes delivered

-- Add validated_dispensers column to deliveries table
ALTER TABLE deliveries 
ADD COLUMN IF NOT EXISTS validated_dispensers JSONB;

-- Add index for potential JSONB queries
CREATE INDEX IF NOT EXISTS idx_deliveries_validated_dispensers ON deliveries USING GIN (validated_dispensers);

-- Comment the column for documentation
COMMENT ON COLUMN deliveries.validated_dispensers IS 'Array of actual dispenser serial codes that were validated and delivered (e.g., ["ABC123", "XYZ456"])';
