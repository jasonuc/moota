import DashboardHeader from '@/components/dashboard-header'
import PlantsList from '@/components/plants-list';
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/plants')({
    component: RouteComponent,
})

function RouteComponent() {
    const seedCount = 10;

    const plants = [
        {
            nickname: "Fernie",
            botanicalName: "Neoregalia Fosteriana",
            hp: 70,
            distance: 50
        },
        {
            nickname: "Leafy",
            botanicalName: "Spathiphyllum Wallisii",
            hp: 60,
            distance: 100
        },
        {
            nickname: "Zoogarte",
            botanicalName: "Ficus elastica",
            hp: 80,
            distance: 150
        },
        {
            nickname: "Sproutlet",
            botanicalName: "Monstera deliciosa",
            hp: 50,
            distance: 200
        },
        {
            nickname: "Spinny",
            botanicalName: "Monstera deliciosa",
            hp: 50,
            distance: 240
        },
        {
            nickname: "Leafy",
            botanicalName: "Spathiphyllum Wallisii",
            hp: 60,
            distance: 100
        },
        {
            nickname: "Zoogarte",
            botanicalName: "Ficus elastica",
            hp: 80,
            distance: 150
        },
        {
            nickname: "Sproutty",
            botanicalName: "Monstera deliciosa",
            hp: 50,
            distance: 200
        }
    ]

    return (
        <div className='flex flex-col space-y-5 pb-10'>
            <DashboardHeader seedCount={seedCount} />

            <h1 className="text-3xl font-heading mb-5">My Plants ({plants.length})</h1>

            <PlantsList plants={plants} maxPlants={plants?.length} />
        </div>
    )
}
