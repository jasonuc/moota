import { Plant } from "./plant";

type User = {
  id: string;
  username: string;
  title: string;
  email: string;
  xp: number;
  level: number;
  createdAt: Date;
  updatedAt: Date;
};

type UserProfile = {
  username: string;
  title: string;
  level: number;
  plantCount: PlantCount;
  seedCount: SeedCount;
  top3AlivePlants: Plant[];
  deceasedPlants: Plant[];
};

type PlantCount = {
  alive: number;
  deceased: number;
};

type SeedCount = {
  Planted: number;
  unused: number;
};

type Stats = {
  plant: PlantCount;
  seed: SeedCount;
};

export type { User, UserProfile, PlantCount, SeedCount, Stats };
