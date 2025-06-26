import { UserProfile } from "@/types/user";
import { ax } from "./index";

export const getUserProfile = async (username: string) =>
  (await ax.get<{ userProfile: UserProfile }>(`/users/${username}/profile`))
    .data.userProfile;

export const getUsernameFromUserId = async (userId: string) =>
  (await ax.get<{ username: string }>(`/users/u/${userId}/username`)).data
    .username;

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

// export const getUser = async (userId: string) =>
//   (await ax.get<{ user: User }>(`/users/u/${userId}/`)).data.user;
