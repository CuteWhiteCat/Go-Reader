-- Add volume-related numbering to chapters
ALTER TABLE chapters ADD COLUMN volume_number INTEGER DEFAULT 1;
ALTER TABLE chapters ADD COLUMN volume_chapter_number INTEGER DEFAULT 0;

-- Backfill existing rows
UPDATE chapters
SET volume_number = 1
WHERE volume_number IS NULL;

UPDATE chapters
SET volume_chapter_number = chapter_number
WHERE volume_chapter_number IS NULL OR volume_chapter_number = 0;
