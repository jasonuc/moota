import { useMutation, useQueryClient } from "@tanstack/react-query";
import { changeUsername, changeEmail, changePassword } from "../api/user";

export const useChangeUsername = () => {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: {
      userId: string;
      currentUsername: string;
      newUsername: string;
    }) => changeUsername(payload.userId, payload.newUsername),
    onSuccess: (_, variables) =>
      queryClient.invalidateQueries({
        queryKey: ["profile", { username: variables.currentUsername }],
        exact: true,
      }),
  });
};

export const useChangeEmail = () =>
  useMutation({
    mutationFn: (payload: { userId: string; newEmail: string }) =>
      changeEmail(payload.userId, payload.newEmail),
  });

export const useChangePassword = () =>
  useMutation({
    mutationFn: (payload: {
      userId: string;
      oldPassword: string;
      newPassword: string;
    }) =>
      changePassword(payload.userId, payload.oldPassword, payload.newPassword),
  });
