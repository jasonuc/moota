-- Create test soils with realistic geographic coordinates
-- Loam soil (medium radius) - Central Park, NYC
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000101', 'loam', 0.55, 0.75, 22.0, 
   ST_GeogFromText('POINT(-73.965355 40.782865)'));

-- Sandy soil (small radius) - Miami Beach
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000102', 'sandy', 0.25, 0.20, 17.5, 
   ST_GeogFromText('POINT(-80.134358 25.792236)'));

-- Silt soil (large radius) - Mississippi Delta
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000103', 'silt', 0.65, 0.55, 30.9, 
   ST_GeogFromText('POINT(-89.254120 29.153291)'));

-- Clay soil (small radius) - Arizona desert
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000104', 'clay', 0.80, 0.65, 17.5, 
   ST_GeogFromText('POINT(-111.928651 33.252440)'));

-- Set of soils close to each other for proximity testing - San Francisco Bay Area
INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000105', 'loam', 0.50, 0.70, 15.0, 
   ST_GeogFromText('POINT(-122.419416 37.774929)')); -- San Francisco

INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000106', 'silt', 0.60, 0.50, 18.0, 
   ST_GeogFromText('POINT(-122.271114 37.804363)')); -- Oakland - ~10km from SF

INSERT INTO soils (id, soil_type, water_retention, nutrient_richness, radius_m, centre) VALUES 
  ('00000000-0000-4000-c000-000000000107', 'clay', 0.75, 0.60, 20.0, 
   ST_GeogFromText('POINT(-122.191291 37.462590)')); -- Palo Alto - ~40km from SF