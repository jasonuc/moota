import { Plant } from "@/types/plant";
import { SeedAvailability, SeedGroup } from "@/types/seed";
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
  await ax.post<{ plant: Plant }>(`/seeds/${seedId}`, {
    latitude,
    longitude,
  });
};

export const requestSeeds = async (userId: string) =>
  (await ax.post<{ seeds: SeedGroup[] }>(`/seeds/u/${userId}/request`)).data
    .seeds;

export const checkWhenUserCanRequestSeeds = async (userId: string) =>
  (await ax.get<SeedAvailability>(`/seeds/u/${userId}/request`)).data;
