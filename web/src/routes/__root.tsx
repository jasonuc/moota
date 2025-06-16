import AuthProvider from "@/services/providers/auth-provider";
import { createRootRoute, Outlet } from "@tanstack/react-router";
import { TanStackRouterDevtools } from "@tanstack/react-router-devtools";

export const Route = createRootRoute({
  component: () => (
    <>
      <AuthProvider>
        <div className="flex font-archivo p-5 md:p-10 min-h-screen w-full">
          <Outlet />
        </div>
        <TanStackRouterDevtools />
      </AuthProvider>
    </>
  ),
});
