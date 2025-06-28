import { AuthContext } from "@/contexts/auth-context";
import { LoginCredentials, RegisterCredentials } from "@/types/auth";
import { User } from "@/types/user";
import { AxiosError } from "axios";
import React, { useCallback, useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { ax } from "../api";

interface AuthProviderProps {
  children: React.ReactNode;
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [{ isLoggedIn, user, error, isInitialized }, setState] = useState<{
    isLoggedIn: boolean;
    user?: User | null;
    error: string | null;
    isInitialized: boolean;
  }>({
    error: null,
    user: null,
    isLoggedIn: false,
    isInitialized: false,
  });
  const navigate = useNavigate();

  const clearError = useCallback(() => {
    setState((prev) => ({
      ...prev,
      error: null,
    }));
  }, []);

  const checkStatus = useCallback(async () => {
    const response = await ax.get<{ user: User | null }>("/whoami");

    setState((prev) => ({
      ...prev,
      isLoggedIn: !!response.data.user,
      user: response.data.user,
      error: null,
      isInitialized: true,
    }));
  }, []);

  useEffect(() => {
    setIsLoading(true);
    checkStatus()
      .catch((err: AxiosError<{ error: string }>) => {
        setState((prev) => ({
          ...prev,
          error: err?.response?.data?.error || "Something went wrong",
          isInitialized: true,
        }));
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [checkStatus]);

  const register = useCallback(
    async (creds: RegisterCredentials) => {
      setIsLoading(true);
      clearError();

      try {
        await ax.post("/auth/register", creds);

        await checkStatus();
        navigate("/home");
      } catch (err) {
        const axiosError = err as AxiosError<{ error: string }>;
        setState((prev) => ({
          ...prev,
          error: axiosError?.response?.data?.error || "Something went wrong",
        }));
      } finally {
        setIsLoading(false);
      }
    },
    [checkStatus, navigate, clearError]
  );

  const login = useCallback(
    async (creds: LoginCredentials) => {
      setIsLoading(true);
      clearError();

      try {
        await ax.post("/auth/login", creds);

        await checkStatus();
        navigate("/home");
      } catch (err) {
        const axiosError = err as AxiosError<{ error: string }>;
        setState((prev) => ({
          ...prev,
          error: axiosError?.response?.data?.error || "Something went wrong",
        }));
      } finally {
        setIsLoading(false);
      }
    },
    [checkStatus, navigate, clearError]
  );

  const refetch = useCallback(async () => {
    setIsLoading(true);
    checkStatus()
      .catch((err: AxiosError<{ error: string }>) => {
        setState((prev) => ({
          ...prev,
          error: err?.response?.data?.error || "Something went wrong",
          isInitialized: true,
        }));
      })
      .finally(() => {
        setIsLoading(false);
      });
  }, [checkStatus]);

  const logout = useCallback(async () => {
    setIsLoading(true);
    clearError();

    try {
      await ax.post("/auth/logout", null, {
        withCredentials: true,
      });

      setState((prev) => ({
        ...prev,
        user: null,
        isLoggedIn: false,
        error: null,
      }));

      navigate("/");
    } catch (err) {
      const axiosError = err as AxiosError<{ error: string }>;
      setState((prev) => ({
        ...prev,
        error: axiosError?.response?.data?.error || "Something went wrong",
      }));
    } finally {
      setIsLoading(false);
    }
  }, [navigate, clearError]);

  const { pathname } = useLocation();

  useEffect(() => {
    const nonProtectedRoutes = ["/", "/login", "/register"];
    if (nonProtectedRoutes.includes(pathname) && !isLoggedIn) return;

    window.addEventListener("auth:logout", logout);
    return () => window.removeEventListener("auth:logout", logout);
  }, [isLoggedIn, pathname, logout]);

  return (
    <AuthContext.Provider
      value={{
        isLoading,
        isInitialized,
        isLoggedIn,
        user,
        error,
        register,
        login,
        refetch,
        logout,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}
