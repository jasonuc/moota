import { CircleMeta } from "./circle";

type SoilType = "loam" | "sandy" | "silt" | "clay";

type SoilMeta = {
  soilType: SoilType;
  waterRetention: number;
  nutrientRichness: number;
};

type Soil = {
  id: string;
  createdAt: Date;
} & SoilMeta &
  CircleMeta;

export type { Soil, SoilType, SoilMeta };
