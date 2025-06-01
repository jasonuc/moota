import { CircleMeta } from "./circle";
import { Soil } from "./soil";

type Plant = {
  id: string;
  nickname: string;
  Hp: number;
  dead: boolean;
  activated: boolean;
  ownerID: string;
  tempers: Tempers;
  timePlanted: Date;
  timeOfDeath: Date;
  lastWateredTime: Date;
  lastActionTime: Date;
} & CircleMeta &
  (Soil | undefined);

type Tempers = {
  woe: number;
  frolic: number;
  malice: number;
  dread: number;
};

export type { Plant };
