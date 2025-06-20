-- Create a single test user
INSERT INTO users (id, username, email, password_hash) VALUES 
  ('00000000-0000-4000-a000-000000000001', 'testuser', 'test@example.com', '\x0123456789ABCDEF');

-- Create test soils
-- Loam soil (medium radius) - Central Park, NYC
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-a000-000000000101', 'loam', 0.55, 0.75, 22.0, 
   ST_GeogFromText('POINT(-73.965355 40.782865)'));

-- Sandy soil (small radius) - Miami Beach
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-a000-000000000102', 'sandy', 0.25, 0.20, 17.5, 
   ST_GeogFromText('POINT(-80.134358 25.792236)'));

-- Silt soil (large radius) - Mississippi Delta
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-a000-000000000103', 'silt', 0.65, 0.55, 30.9, 
   ST_GeogFromText('POINT(-89.254120 29.153291)'));

-- Clay soil (small radius) - Arizona desert
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-a000-000000000104', 'clay', 0.80, 0.65, 17.5, 
   ST_GeogFromText('POINT(-111.928651 33.252440)'));

-- Plants in loam soil (Central Park)
-- Center plant
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000201', 'Oak Tree', 100.0, false, 
  '00000000-0000-4000-a000-000000000001', NOW(),
  ST_GeogFromText('POINT(-73.965355 40.782865)'), 3.0, 
  '00000000-0000-4000-a000-000000000101', 'loam', 'Quercus alba', 
  3, 2, 4, 1
);

-- Plant 8m north of center (still within soil radius)
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000202', 'Pine', 95.0, false, 
  '00000000-0000-4000-a000-000000000001', NOW(),
  ST_GeogFromText('POINT(-73.965355 40.782937)'), 3.0, 
  '00000000-0000-4000-a000-000000000101', 'loam', 'Pinus strobus', 
  2, 5, 1, 3
);

-- Plant 8m east of center (still within soil radius)
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000203', 'Maple', 90.0, false, 
  '00000000-0000-4000-a000-000000000001', NOW(),
  ST_GeogFromText('POINT(-73.965265 40.782865)'), 3.0, 
  '00000000-0000-4000-a000-000000000101', 'silt', 'Acer rubrum', 
  1, 4, 3, 5
);

-- Plants in sandy soil (Miami Beach)
-- Center plant
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000204', 'Palm Tree', 85.0, false, 
  '00000000-0000-4000-a000-000000000001', NOW(),
  ST_GeogFromText('POINT(-80.134358 25.792236)'), 3.0, 
  '00000000-0000-4000-a000-000000000102', 'sandy', 'Cocos nucifera', 
  5, 3, 2, 1
);

-- Plant 7m south of center (still within soil radius)
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000205', 'Cactus', 75.0, false, 
  '00000000-0000-4000-a000-000000000001', NOW(),
  ST_GeogFromText('POINT(-80.134358 25.792173)'), 3, 
  '00000000-0000-4000-a000-000000000102', 'sandy', 'Opuntia', 
  4, 1, 5, 2
);

-- Dead plant in clay soil
INSERT INTO plants (
  id, nickname, hp, dead, owner_id, time_planted, centre, radius_m, 
  soil_id, optimal_soil, botanical_name, woe, frolic, dread, malice
) VALUES (
  '00000000-0000-4000-a000-000000000206', 'Wilted Flower', 0.0, true, 
  '00000000-0000-4000-a000-000000000001', NOW() - INTERVAL '10 days', 
  ST_GeogFromText('POINT(-111.928651 33.252440)'), 3.0, 
  '00000000-0000-4000-a000-000000000104', 'loam', 'Rosa rubiginosa', 
  2, 1, 5, 3
);
