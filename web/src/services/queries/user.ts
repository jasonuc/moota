import { useQuery } from "@tanstack/react-query";
import { getUsernameFromUserId, getUserProfile } from "../api/user";
import { useAuth } from "@/hooks/use-auth";

export const useGetUserProfile = (username?: string) =>
  useQuery({
    queryKey: ["profile", { username }],
    queryFn: () => getUserProfile(username!),
    enabled: !!username,
  });

export const useGetCurrentUserProfile = () => {
  const { user } = useAuth();
  return useQuery({
    queryKey: ["current-user-profile", { username: user?.username }],
    queryFn: () => getUserProfile(user!.username),
    enabled: !!user,
  });
};

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
