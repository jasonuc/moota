import { LogoWithText } from '@/components/logo'
import { createFileRoute, Link, Outlet } from '@tanstack/react-router'

export const Route = createFileRoute('/_auth')({
    component: RouteComponent,
})

function RouteComponent() {
    return (
        <div className='w-full flex flex-col gap-y-10'>
            <div className=''>
                <Link to='/'>
                    <LogoWithText />
                </Link>
            </div>

            <div className='flex h-full md:items-center justify-center'>
                <Outlet />
            </div>
        </div>
    )
}
