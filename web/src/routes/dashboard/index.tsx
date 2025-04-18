import DashboardHeader from '@/components/dashboard-header'
import MyNearbyPlants from '@/components/my-nearby-plants';
import { Button } from '@/components/ui/button';
import { createFileRoute } from '@tanstack/react-router'
import { BeanIcon, SproutIcon } from 'lucide-react';

export const Route = createFileRoute('/dashboard/')({
    component: RouteComponent,
})

function RouteComponent() {
    const seedCount = 10;

    return (
        <div className='flex flex-col space-y-5'>
            <DashboardHeader seedCount={seedCount} />

            <MyNearbyPlants />

            <div className="fixed left-0 bottom-0 w-full flex gap-x-5 p-5 md:p-10">
                <Button className='hover:cursor-pointer grow'>
                    My Plants
                    <SproutIcon className="ml-2" />
                </Button>

                <Button className='hover:cursor-pointer grow'>
                    Plant Seed
                    <BeanIcon className="ml-2" />
                </Button>
            </div>
        </div>
    )
}