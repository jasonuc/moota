import { Coordinates } from "./coordinates";
import { LevelMeta } from "./level";
import { Soil } from "./soil";

type Plant = {
  id: string;
  nickname: string;
  botanicalName: string;
  hp: number;
  dead: boolean;
  ownerID: string;
  tempers: Tempers;
  timePlanted: string;
  timeOfDeath: string;
  lastWateredTime: string;
  lastActionTime: string;
  centre: Coordinates;
  soil: Soil;
} & LevelMeta;

type Tempers = {
  woe: number;
  frolic: number;
  malice: number;
  dread: number;
};

type PlantWithDistanceMFromUser = Plant & {
  distanceM?: number;
};

export type { Plant, PlantWithDistanceMFromUser };
