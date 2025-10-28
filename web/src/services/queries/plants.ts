import { keepPreviousData, useQuery } from "@tanstack/react-query";
import {
  getUserDeceasedPlants,
  getPlant,
  getUserNearbyPlants,
  getUserPlants,
} from "../api/plants";
import { AxiosError } from "axios";
import { SERVER_404_MESSAGE } from "@/lib/constants";

export const useGetUserNearbyPlants = (
  userId?: string,
  lat?: number | null,
  lon?: number | null,
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
  lon?: number | null,
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
    retry: (failureCount, error) => {
      const err = error as AxiosError<{ error: string }>;
      return (
        failureCount < 3 && err.response?.data.error !== SERVER_404_MESSAGE
      );
    },
  });

export const useGetUserDeceasedPlants = (userId?: string) =>
  useQuery({
    queryKey: ["plants", { userId, deceased: true }],
    queryFn: () => getUserDeceasedPlants(userId!),
    enabled: !!userId,
  });
