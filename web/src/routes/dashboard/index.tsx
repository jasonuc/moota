import HomeHeader from '@/components/home-header'
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
        <div className='flex flex-col space-y-6'>
            <HomeHeader seedCount={seedCount} />

            <MyNearbyPlants />

            <div className="w-full grid grid-cols-2 gap-3">
                <Button className='hover:cursor-pointer'>
                    My Plants
                    <SproutIcon className="ml-2" />
                </Button>

                <Button className='hover:cursor-pointer'>
                    Plant Seed
                    <BeanIcon className="ml-2" />
                </Button>
            </div>
        </div>
    )
}