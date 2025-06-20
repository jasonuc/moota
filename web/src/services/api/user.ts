import { ax } from "./index";
import { UserProfile } from "@/types/user";

export const getUserProfile = async (username: string) =>
  (await ax.get<{ userProfile: UserProfile }>(`/users/${username}/profile`))
    .data.userProfile;

export const changeUsername = async (userId: string, newUsername: string) =>
  await ax.patch(`/auth/u/${userId}/change-username`, {
    newUsername,
  });

export const changeEmail = async (userId: string, newEmail: string) =>
  await ax.patch(`/auth/u/${userId}/change-email`, {
    newEmail,
  });

export const changePassword = async (
  userId: string,
  oldPassword: string,
  newPassword: string
) =>
  await ax.patch(`/auth/u/${userId}/change-password`, {
    oldPassword,
    newPassword,
  });
