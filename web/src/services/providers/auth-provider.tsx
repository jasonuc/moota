import { AuthContext } from "@/contexts/auth-context";
import { LoginCredentials, RegisterCredentials } from "@/types/auth";
import { User } from "@/types/user";
import { useNavigate } from "react-router";
import { AxiosError } from "axios";
import { ax } from "../api";
import React, { useCallback, useEffect, useState } from "react";

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
    checkStatus().catch((err: AxiosError<{ error: string }>) => {
      setState((prev) => ({
        ...prev,
        error: err?.response?.data?.error || "Something went wrong",
      }));
    });
    setIsLoading(false);
  }, [checkStatus]);

  const register = useCallback(
    async (creds: RegisterCredentials) => {
      setIsLoading(true);
      ax.post("/auth/register", creds, {
        withCredentials: true,
      })
        .then(() => {
          checkStatus();
          navigate("/home");
        })
        .catch((err: AxiosError<{ error: string }>) => {
          setState((prev) => ({
            ...prev,
            error: err?.response?.data?.error || "Something went wrong",
          }));
        });
      setIsLoading(false);
    },
    [checkStatus, navigate]
  );

  const login = useCallback(
    async (creds: LoginCredentials) => {
      setIsLoading(true);
      ax.post("/auth/login", creds, {
        withCredentials: true,
      })
        .then(() => {
          checkStatus();
          navigate("/home");
        })
        .catch((err: AxiosError<{ error: string }>) => {
          setState((prev) => ({
            ...prev,
            error: err?.response?.data?.error || "Something went wrong",
          }));
        });
      setIsLoading(false);
    },
    [checkStatus, navigate]
  );

  const logout = useCallback(async () => {
    ax.post("/auth/logout", null, {
      withCredentials: true,
    })
      .then(() => {
        checkStatus();
        navigate("/");
      })
      .catch((err: AxiosError<{ error: string }>) => {
        setState((prev) => ({
          ...prev,
          error: err?.response?.data?.error || "Something went wrong",
        }));
      });
    setIsLoading(false);
  }, [checkStatus, navigate]);

  return (
    <AuthContext
      value={{ isLoading, isLoggedIn, user, error, login, register, logout }}
    >
      {children}
    </AuthContext>
  );
}
