CREATE TABLE IF NOT EXISTS plants (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid (),
	nickname VARCHAR(255) NOT NULL,
	hp REAL NOT NULL,
	dead BOOLEAN NOT NULL DEFAULT FALSE,
	owner_id UUID REFERENCES users (id) ON DELETE CASCADE,
	time_planted TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_watered_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_action_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
	last_refreshed_at TIMESTAMPTZ,
	grace_period_ends_at TIMESTAMPTZ,
	time_of_death TIMESTAMPTZ,
	centre GEOGRAPHY (POINT) NOT NULL,
	radius_m DOUBLE PRECISION NOT NULL,
	soil_id UUID REFERENCES soils (id) ON DELETE CASCADE,
	optimal_soil VARCHAR(30) NOT NULL,
	botanical_name TEXT NOT NULL,
	level SMALLINT NOT NULL DEFAULT 1,
	xp INTEGER NOT NULL DEFAULT 0,
	woe REAL NOT NULL CHECK (woe >= 1 AND woe <= 5),
	frolic REAL NOT NULL CHECK (frolic >= 1 AND frolic <= 5),
	dread REAL NOT NULL CHECK (dread >= 1 AND dread <= 5),
	malice REAL NOT NULL CHECK (malice >= 1 AND malice <= 5)
);