import { Outlet } from "react-router";

export default function BaseLayout() {
  return (
    <div className="flex font-archivo p-5 md:p-10 min-h-screen w-full">
      <Outlet />
    </div>
  );
}
