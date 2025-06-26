import { useGetCurrentUserProfile } from "@/services/queries/user";

export function useCurrentUserProfile() {
  const { data: currentUserProfile } = useGetCurrentUserProfile();
  return currentUserProfile;
}
