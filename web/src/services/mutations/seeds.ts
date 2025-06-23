import { useMutation, useQueryClient } from "@tanstack/react-query";
import { plantSeed, requestSeeds } from "../api/seeds";
import { useAuth } from "@/hooks/use-auth";

export const usePlantSeed = () => {
  const { user } = useAuth();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { seedId: string; lat: number; lon: number }) =>
      plantSeed(payload.seedId, payload.lat, payload.lon),
    onSuccess: () =>
      queryClient.invalidateQueries({
        queryKey: ["seeds", { userId: user?.id }],
        exact: true,
      }),
  });
};

export const useRequestSeeds = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (userId: string) => requestSeeds(userId),
    onSuccess: (_, userId) =>
      queryClient.invalidateQueries({
        queryKey: ["seeds", { userId }],
        exact: true,
      }),
  });
};
