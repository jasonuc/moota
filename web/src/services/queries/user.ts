import { useQuery } from "@tanstack/react-query";
import { getUsernameFromUserId, getUserProfile } from "../api/user";

export const useGetUserProfile = (username?: string) =>
  useQuery({
    queryKey: ["profile", { username }],
    queryFn: () => getUserProfile(username!),
    enabled: !!username,
  });

export const useGetCurrentUserProfile = (username?: string) =>
  useQuery({
    queryKey: ["current-user-profile", { username }],
    queryFn: () => getUserProfile(username!),
  });

export const useGetUsernameFromUserId = (userId?: string) =>
  useQuery({
    queryKey: ["user-username", { userId }],
    queryFn: () => getUsernameFromUserId(userId!),
    enabled: !!userId,
  });

// export const useGetUser = (userId?: string) =>
//   useQuery({
//     queryKey: ["user", { userId }],
//     queryFn: () => getUser(userId!),
//     enabled: !!userId,
//   });
