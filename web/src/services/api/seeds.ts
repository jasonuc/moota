import { SeedGroup } from "@/types/seed";
import { ax } from "./index";

export const getUserSeeds = async (userId: string) => {
  return (await ax.get<{ seeds: SeedGroup[] }>(`/seeds/u/${userId}`)).data
    .seeds;
};

export const plantSeed = async (
  seedId: string,
  latitude: number,
  longitude: number
) => {
  await ax.post(`/seeds/${seedId}`, {
    latitude: latitude,
    longitude: longitude,
  });
};
