-- ============================================
-- Add branch type
-- ============================================

ALTER TABLE branches
ADD COLUMN type VARCHAR(20) NOT NULL DEFAULT 'SUPERMARKET';
