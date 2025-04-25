import Logo from '@/components/logo'
import { createFileRoute, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <div className='w-full flex flex-col gap-y-10'>
            <div className=''>
                <Logo />
            </div>

            <Outlet />
        </div>
    )
}
