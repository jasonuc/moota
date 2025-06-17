import { LoginCredentials, RegisterCredentials } from "@/types/auth";
import { User } from "@/types/user";
import { createContext } from "react";

export const AuthContext = createContext<{
  user?: User | null;
  isLoading: boolean;
  isLoggedIn: boolean;
  isInitialized: boolean;
  error?: string | null;
  login: (creds: LoginCredentials) => Promise<void>;
  logout: () => void;
  register: (creds: RegisterCredentials) => Promise<void>;
}>({
  isLoading: true,
  isLoggedIn: false,
  isInitialized: false,
  login: async () => {},
  logout: () => {},
  register: async () => {},
});
