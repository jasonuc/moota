import { LoginCredentials, RegisterCredentials } from "@/types/auth";
import { User } from "@/types/user";
import { createContext } from "react";

export const AuthContext = createContext<{
  user?: User | null;
  isLoading: boolean;
  isLoggedIn: boolean;
  isInitialized: boolean;
  error?: string | null;
  register: (creds: RegisterCredentials) => Promise<void>;
  login: (creds: LoginCredentials) => Promise<void>;
  refetch: () => Promise<void>;
  logout: () => void;
}>({
  isLoading: true,
  isLoggedIn: false,
  isInitialized: false,
  register: async () => {},
  login: async () => {},
  refetch: async () => {},
  logout: () => {},
});
