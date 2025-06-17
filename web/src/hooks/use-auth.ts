import { AuthContext } from "@/contexts/auth-context";
import { useContext } from "react";

export const useAuth = () => useContext(AuthContext);
