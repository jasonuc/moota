import axios from "axios";
import { useState } from "react";

// Dicebear is used most images in dynamic images in this application like plants or user profile photos
export function useIsDicebearOnline() {
  const [isOnline, setIsOnline] = useState<boolean>();
  axios
    .get("https://api.dicebear.com")
    .then(() => setIsOnline(true))
    .catch(() => setIsOnline(false));

  return isOnline;
}
