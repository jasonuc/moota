import { PlantWithDistanceMFromUser } from "@/types/plant";
import { ax } from "./index";

export const getUserNearbyPlants = async (
  userId: string,
  lat: number,
  lon: number
) =>
  (
    await ax.get<{ plants: PlantWithDistanceMFromUser[] }>(
      `/plants/u/${userId}`,
      {
        params: { lat: lat, lon: lon },
      }
    )
  ).data.plants.slice(0, 4);

export const getAllUserPlants = async (
  userId: string,
  lat: number,
  lon: number
) =>
  (
    await ax.get<{ plants: PlantWithDistanceMFromUser[] }>(
      `/plants/u/${userId}`,
      {
        params: { lat: lat, lon: lon },
      }
    )
  ).data.plants;
