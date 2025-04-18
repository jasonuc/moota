import HomeHeader from '@/components/home-header'
import MyNearbyPlants from '@/components/my-nearby-plants';
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/dashboard/')({
    component: RouteComponent,
})

function RouteComponent() {
    const seedCount = 10;

    return (
        <div className='flex flex-col space-y-2'>
            <HomeHeader seedCount={seedCount} />

            <MyNearbyPlants />
        </div>
    )
}
