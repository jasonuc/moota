import { useMutation, useQueryClient } from "@tanstack/react-query";
import { changePlantNickname, killPlant, waterPlant } from "../api/plants";

export const useKillPlant = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (plantId: string) => killPlant(plantId),
    onSuccess: (_, plantId) =>
      queryClient.removeQueries({
        queryKey: ["plant", plantId],
        exact: true,
      }),
  });
};

export const useChangePlantNickname = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { plantId: string; newNickname: string }) =>
      changePlantNickname(payload.plantId, payload.newNickname),
    onSuccess: (_, variables) =>
      queryClient.removeQueries({
        queryKey: ["plant", variables.plantId],
        exact: true,
      }),
  });
};

export const useWaterPlant = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      plantId: string;
      latitude: number;
      longitude: number;
    }) => waterPlant(payload.plantId, payload.latitude, payload.longitude),
    onSuccess: (_, variables) =>
      queryClient.refetchQueries({
        queryKey: ["plant", variables.plantId],
        exact: true,
      }),
  });
};
