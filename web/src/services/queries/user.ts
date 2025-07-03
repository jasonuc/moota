import { useQuery } from "@tanstack/react-query";
import { getUsernameFromUserId, getUserProfile } from "../api/user";
import { useAuth } from "@/hooks/use-auth";
import { AxiosError } from "axios";
import { SERVER_404_MESSAGE } from "@/lib/constants";

export const useGetUserProfile = (username?: string) =>
  useQuery({
    queryKey: ["profile", { username }],
    queryFn: () => getUserProfile(username!),
    enabled: !!username,
    retry: (failureCount, error) => {
      const err = error as AxiosError<{ error: string }>;
      return (
        failureCount < 3 && err.response?.data.error !== SERVER_404_MESSAGE
      );
    },
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
