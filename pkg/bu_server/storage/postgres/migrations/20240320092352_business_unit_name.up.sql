CREATE EXTENSION IF NOT EXISTS pg_trgm;

ALTER TABLE business_unit ADD COLUMN IF NOT EXISTS "name" TEXT;
ALTER TABLE trade_document ADD COLUMN IF NOT EXISTS "doc_reference" TEXT DEFAULT '' NOT NULL;

UPDATE business_unit SET "name" = "business_unit"->>'name';
ALTER TABLE business_unit ALTER COLUMN "name" SET NOT NULL;

CREATE INDEX IF NOT EXISTS business_unit_name_idx ON business_unit USING GIN(name gin_trgm_ops);
CREATE INDEX IF NOT EXISTS trade_document_doc_reference_idx ON trade_document USING GIN(doc_reference gin_trgm_ops);
