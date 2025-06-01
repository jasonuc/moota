import { Coordinates } from "./coordinates";

type CircleMeta = {
  c: Coordinates;
  radiusM: number;
};

interface Circle {
  Centre: () => Coordinates;
  RadiusM: () => number;
  ContainsPoint: (arg0: Coordinates) => boolean;
  OverlapsWith: (arg0: Circle) => boolean;
}

export type { CircleMeta, Circle };
