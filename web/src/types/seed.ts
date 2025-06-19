import { SoilType } from "./soil";

type SeedMeta = {
  botanicalName: string;
  optimalSoil: SoilType;
};

type Seed = {
  id: string;
  hp: number;
  planted: boolean;
  ownerID: string;
  createdAt: string;
} & SeedMeta;

type SeedGroup = {
  botanicalName: string;
  count: number;
  seeds: Seed[];
};

type SeedAvailability = {
  timeAvailable: string;
  availableNow: boolean;
};

export type { Seed, SeedAvailability, SeedGroup, SeedMeta };
