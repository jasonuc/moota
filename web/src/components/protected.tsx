import { useAuth } from "@/hooks/use-auth";
import React from "react";
import { Navigate } from "react-router";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

export default function ProtectedRoute({ children }: ProtectedRouteProps) {
  const { isLoggedIn, isLoading } = useAuth();

  if (isLoading) return null;
  if (!isLoggedIn) return <Navigate to="/" />;

  return children;
}
