ALTER TABLE diagnostic_centres
ALTER COLUMN doctors TYPE TEXT[] USING doctors::TEXT[],
ALTER COLUMN available_tests TYPE TEXT[] USING available_tests::TEXT[],
ALTER COLUMN doctors SET DEFAULT ARRAY['Male']::TEXT[],
ALTER COLUMN available_tests SET DEFAULT ARRAY['CT_SCAN']::TEXT[];

