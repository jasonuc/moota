import { keepPreviousData, useQuery } from "@tanstack/react-query";
import {
  getUserDeceasedPlants,
  getPlant,
  getUserNearbyPlants,
  getUserPlants,
} from "../api/plants";

export const useGetUserNearbyPlants = (
  userId?: string,
  lat?: number | null,
  lon?: number | null
) =>
  useQuery({
    queryKey: ["plants", { userId, lat, lon, deceased: false }],
    queryFn: () => getUserNearbyPlants(userId!, lat ?? 0, lon ?? 0),
    enabled: () => !!userId,
    placeholderData: keepPreviousData,
  });

export const useGetUserPlants = (
  userId?: string,
  lat?: number | null,
  lon?: number | null
) =>
  useQuery({
    queryKey: ["plants", { userId, lat, lon, deceased: false }],
    queryFn: () => getUserPlants(userId!, lat ?? 0, lon ?? 0),
    enabled: () => !!userId,
    placeholderData: keepPreviousData,
  });

export const useGetPlant = (plantId?: string) =>
  useQuery({
    queryKey: ["plant", { plantId: plantId }],
    queryFn: () => getPlant(plantId!),
    enabled: !!plantId,
  });

export const useGetUserDeceasedPlants = (userId?: string) =>
  useQuery({
    queryKey: ["plants", { userId, deceased: true }],
    queryFn: () => getUserDeceasedPlants(userId!),
    enabled: !!userId,
  });
