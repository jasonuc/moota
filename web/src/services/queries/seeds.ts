import { useQuery } from "@tanstack/react-query";
import { checkWhenUserCanRequestSeeds, getUserSeeds } from "../api/seeds";

export const useGetUserSeeds = (userId?: string) =>
  useQuery({
    queryKey: ["seeds", { userId }],
    queryFn: () => getUserSeeds(userId!),
    enabled: !!userId,
  });

export const useCheckWhenUserCanRequestSeeds = (userId?: string) =>
  useQuery({
    queryKey: ["seed-request-status", { userId: userId! }],
    queryFn: () => checkWhenUserCanRequestSeeds(userId!),
    enabled: !!userId,
  });
