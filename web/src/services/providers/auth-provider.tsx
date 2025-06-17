import { AuthContext } from "@/contexts/auth-context";
import { LoginCredentials, RegisterCredentials } from "@/types/auth";
import { User } from "@/types/user";
import { AxiosError } from "axios";
import React, { useCallback, useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { ax } from "../api";

interface AuthProviderProps {
  children: React.ReactNode;
}

export default function AuthProvider({ children }: AuthProviderProps) {
  const [isLoading, setIsLoading] = useState(false);
  const [{ isLoggedIn, user, error }, setState] = useState<{
    isLoggedIn: boolean;
    user?: User | null;
    error: string | null;
  }>({
    error: null,
    user: null,
    isLoggedIn: false,
  });
  const navigate = useNavigate();

  const clearError = useCallback(() => {
    setState((prev) => ({
      ...prev,
      error: null,
    }));
  }, []);

  const checkStatus = useCallback(async () => {
    const response = await ax.get<{ user: User | null }>("/whoami", {
      withCredentials: true,
    });

    setState((prev) => ({
      ...prev,
      isLoggedIn: !!response.data.user,
      user: response.data.user,
      error: null,
    }));
  }, []);

  useEffect(() => {
    setIsLoading(true);
    checkStatus()
      .catch((err: AxiosError<{ error: string }>) => {
        setState((prev) => ({
          ...prev,
          error: err?.response?.data?.error || "Something went wrong",
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
        await ax.post("/auth/register", creds, {
          withCredentials: true,
        });

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
        await ax.post("/auth/login", creds, {
          withCredentials: true,
        });

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

  return (
    <AuthContext.Provider
      value={{ isLoading, isLoggedIn, user, error, login, register, logout }}
    >
      {children}
    </AuthContext.Provider>
  );
}
