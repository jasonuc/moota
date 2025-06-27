import { useMutation, useQueryClient } from "@tanstack/react-query";
import { changeUsername, changeEmail, changePassword } from "../api/user";
import { useAuth } from "@/hooks/use-auth";

export const useChangeUsername = () => {
  const { refetch } = useAuth();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      userId: string;
      currentUsername: string;
      newUsername: string;
    }) => changeUsername(payload.userId, payload.newUsername),
    onSuccess: () => {
      queryClient.refetchQueries({
        queryKey: ["current-user-profile"],
        exact: false,
      });
      refetch();
    },
  });
};

export const useChangeEmail = () => {
  const { refetch } = useAuth();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: { userId: string; newEmail: string }) =>
      changeEmail(payload.userId, payload.newEmail),
    onSuccess: (_, variables) => {
      queryClient.refetchQueries({
        queryKey: ["user", { userId: variables.userId }],
      });
      refetch();
    },
  });
};

export const useChangePassword = () => {
  const { refetch } = useAuth();
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      userId: string;
      oldPassword: string;
      newPassword: string;
    }) =>
      changePassword(payload.userId, payload.oldPassword, payload.newPassword),
    onSuccess: (_, variables) => {
      queryClient.refetchQueries({
        queryKey: ["user", { userId: variables.userId }],
      });
      refetch();
    },
  });
};
