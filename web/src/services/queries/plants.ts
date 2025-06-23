import { useQuery } from "@tanstack/react-query";
import {
  getAllUserDeceasedPlants,
  getAllUserPlants,
  getPlant,
  getUserNearbyPlants,
} from "../api/plants";

export const useGetUserNearbyPlants = (
  userId?: string,
  lat?: number | null,
  lon?: number | null
) =>
  useQuery({
    queryKey: ["plants", { userId, lat, lon, deceased: false }],
    queryFn: () => getUserNearbyPlants(userId!, lat!, lon!),
    enabled: () => !(!userId || !lat || !lon),
  });

export const useGetAllUserPlants = (
  userId?: string,
  lat?: number | null,
  lon?: number | null
) =>
  useQuery({
    queryKey: ["plants", { userId, lat, lon, deceased: false }],
    queryFn: () => getAllUserPlants(userId!, lat!, lon!),
    enabled: () => !(!userId || !lat || !lon),
  });

export const useGetPlant = (plantId?: string) =>
  useQuery({
    queryKey: ["plant", plantId],
    queryFn: () => getPlant(plantId!),
    enabled: !!plantId,
  });

export const useGetAllUserDeceasedPlants = (userId?: string) =>
  useQuery({
    queryKey: ["plants", { userId, deceased: true }],
    queryFn: () => getAllUserDeceasedPlants(userId!),
    enabled: !!userId,
  });
