import { ax } from "./index";
import { UserProfile } from "@/types/user";

export const getUserProfile = async (username: string) =>
  (await ax.get<{ userProfile: UserProfile }>(`/users/${username}/profile`))
    .data.userProfile;
