import { LogoWithText } from "@/components/logo";
import { useAuth } from "@/hooks/use-auth";
import { Link, Navigate, Outlet } from "react-router";

function AuthLayout() {
  const { isLoggedIn } = useAuth();

  if (!isLoggedIn) return <Navigate to={"/home"} />;

  return (
    <div className="w-full flex flex-col gap-y-10">
      <div className="">
        <Link to="/">
          <LogoWithText />
        </Link>
      </div>

      <div className="flex h-full md:items-center justify-center">
        <Outlet />
      </div>
    </div>
  );
}

export default AuthLayout;
