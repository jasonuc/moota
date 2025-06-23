import { useGetUserProfile } from "@/services/queries/user";
import { useAuth } from "./use-auth";

export function useCurrentUserProfile() {
  const { user } = useAuth();
  const { data: currentUserProfile } = useGetUserProfile(user?.username);

  return currentUserProfile;
}
