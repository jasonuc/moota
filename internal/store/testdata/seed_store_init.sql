-- Create test user
INSERT INTO users (id, username, email, password_hash) VALUES 
  ('00000000-0000-4000-a000-000000000001', 'testuser', 'test@example.com', '\x0123456789ABCDEF');

-- Create test seeds
-- Unplanted seeds for testuser
INSERT INTO seeds (id, owner_id, hp, planted, optimal_soil, botanical_name, created_at) VALUES 
  ('00000000-0000-4000-b000-000000000001', '00000000-0000-4000-a000-000000000001', 50.0, false, 'loam', 'Quercus alba', NOW() - INTERVAL '10 days'),
  ('00000000-0000-4000-b000-000000000002', '00000000-0000-4000-a000-000000000001', 50.0, false, 'sandy', 'Cocos nucifera', NOW() - INTERVAL '8 days'),
  ('00000000-0000-4000-b000-000000000003', '00000000-0000-4000-a000-000000000001', 50.0, false, 'silt', 'Acer rubrum', NOW() - INTERVAL '5 days');

-- Already planted seed
INSERT INTO seeds (id, owner_id, hp, planted, optimal_soil, botanical_name, created_at) VALUES 
  ('00000000-0000-4000-b000-000000000004', '00000000-0000-4000-a000-000000000001', 50.0, true, 'clay', 'Rosa rubiginosa', NOW() - INTERVAL '15 days');

-- Seed owned by another user
INSERT INTO users (id, username, email, password_hash) VALUES 
  ('00000000-0000-4000-a000-000000000002', 'testuser2', 'test2@example.com', '\x0123456789ABCDEF');

INSERT INTO seeds (id, owner_id, hp, planted, optimal_soil, botanical_name, created_at) VALUES 
  ('00000000-0000-4000-b000-000000000005', '00000000-0000-4000-a000-000000000002', 50.0, false, 'loam', 'Pinus strobus', NOW() - INTERVAL '3 days');