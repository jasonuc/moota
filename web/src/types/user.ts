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
};

type PlantCount = {
  alive: number;
  deceased: number;
};

type SeedCount = {
  Planted: number;
  unused: number;
};

export type { User, UserProfile };
