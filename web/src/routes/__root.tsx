import { createRootRoute, Outlet } from '@tanstack/react-router'
import { TanStackRouterDevtools } from '@tanstack/react-router-devtools'

export const Route = createRootRoute({
    component: () => (
        <>
            <div className="font-archivo p-5 bg-background min-h-screen w-full">
                <Outlet />
            </div>
            <TanStackRouterDevtools />
        </>
    ),
})