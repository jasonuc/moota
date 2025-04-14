CREATE TABLE IF NOT EXISTS soils (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	soil_type VARCHAR(30) NOT NULL,
	water_retention REAL NOT NULL,
	nutrient_richness REAL NOT NULL,
	radius_m DOUBLE PRECISION NOT NULL,
	centre GEOGRAPHY (POINT) NOT NULL,
	created_at TIMESTAMPTZ DEFAULT NOW()
);