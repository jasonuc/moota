import DashboardHeader from '@/components/dashboard-header'
import MyNearbyPlants from '@/components/my-nearby-plants';
import { Button } from '@/components/ui/button';
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { BeanIcon, SproutIcon } from 'lucide-react';

export const Route = createFileRoute('/dashboard/')({
    component: RouteComponent,
})

function RouteComponent() {
    const seedCount = 10;
    const navigate = useNavigate({ from: '/dashboard' })

    return (
        <div className='flex flex-col space-y-5'>
            <DashboardHeader seedCount={seedCount} />

            <MyNearbyPlants />

            <div className="fixed left-0 bottom-0 w-full flex gap-x-5 p-5 md:p-10">
                <Button
                    onClick={() => navigate({ to: "/dashboard" })}
                    className='hover:cursor-pointer grow'>
                    Plants
                    <SproutIcon className="ml-2" />
                </Button>

                <Button
                    onClick={() => navigate({ to: "/dashboard" })}
                    className='hover:cursor-pointer grow'>
                    Seeds
                    <BeanIcon className="ml-2" />
                </Button>
            </div>
        </div>
    )
}