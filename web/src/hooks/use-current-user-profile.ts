import { useEffect, useState } from "react";
import { useAuth } from "./use-auth";
import { UserProfile } from "@/types/user";
import { getUserProfile } from "@/services/api/user";

export default function useCurrentUserProfile() {
  const { user } = useAuth();
  const [currentUserProfile, setCurrentUserProfile] = useState<UserProfile>();

  useEffect(() => {
    getUserProfile(user!.username).then(setCurrentUserProfile);
  }, [user]);

  return currentUserProfile;
}
