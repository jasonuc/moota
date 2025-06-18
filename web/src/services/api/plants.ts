import { Plant, PlantWithDistanceMFromUser } from "@/types/plant";
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

export const getPlant = async (plantId: string) =>
  (await ax.get<{ plant: Plant }>(`/plants/${plantId}`)).data.plant;

export const killPlant = async (plantId: string) =>
  (await ax.post(`/plants/${plantId}/kill`)).status;

export const waterPlant = async (
  plantId: string,
  latitude: number,
  longitude: number
) =>
  (
    await ax.post<{ plant: Plant }>(`plants/${plantId}/action`, {
      latitude,
      longitude,
      action: 1,
    })
  ).data;

export const getAllUserUnactivatedPlants = async (
  userId: string,
  lat: number,
  lon: number
) =>
  (
    await ax.get<{ plants: PlantWithDistanceMFromUser[] }>(
      `/plants/u/${userId}`,
      {
        params: { lat: lat, lon: lon, includeInactive: true },
      }
    )
  ).data.plants.filter((plant) => !plant.activated);

export const getUserUnactivatedPlantsCount = async (userId: string) =>
  (await getAllUserUnactivatedPlants(userId, 0, 0)).length;

export const activatePlant = async (plantId: string) =>
  (await ax.post<{ plant: Plant }>(`/plants/${plantId}/activate`)).data.plant;
