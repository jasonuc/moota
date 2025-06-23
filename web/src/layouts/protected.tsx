import { useAuth } from "@/hooks/use-auth";
import { Navigate, Outlet } from "react-router";

export default function ProtectedLayout() {
  const { isLoggedIn, isInitialized } = useAuth();

  if (!isInitialized) return null;
  if (!isLoggedIn) return <Navigate to="/" />;

  return <Outlet />;
}
