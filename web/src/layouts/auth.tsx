import { LogoWithText } from "@/components/logo";
import { Link, Outlet } from "react-router";

function AuthLayout() {
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
