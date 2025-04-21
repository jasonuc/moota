import Header from '@/components/header'
import PlantsList from '@/components/plants-list';
import { Button } from '@/components/ui/button';
import { createFileRoute, useNavigate } from '@tanstack/react-router'
import { BeanIcon, SproutIcon } from 'lucide-react';

export const Route = createFileRoute('/dashboard')({
  component: RouteComponent,
})

function RouteComponent() {
  const seedCount = 10;
  const navigate = useNavigate({ from: '/dashboard' })

  const nearbyPlants = [
    {
      id: "1",
      nickname: "Fernie",
      botanicalName: "Neoregalia Fosteriana",
      hp: 70,
      distance: 50
    },
    {
      id: "2",
      nickname: "Leafy",
      botanicalName: "Spathiphyllum Wallisii",
      hp: 60,
      distance: 100
    },
    {
      id: "3",
      nickname: "Zoogarte",
      botanicalName: "Ficus elastica",
      hp: 80,
      distance: 150
    },
    {
      id: "4",
      nickname: "Sproutlet",
      botanicalName: "Monstera deliciosa",
      hp: 50,
      distance: 200
    },
  ]

  return (
    <div className='flex flex-col space-y-5 grow'>
      <Header seedCount={seedCount} />

      <h1 className="text-3xl font-heading mb-5">My Nearby Plants</h1>

      <PlantsList plants={nearbyPlants} />

      <div className='pt-20' />

      <div className="fixed left-0 bottom-0 w-full flex gap-x-5 p-5 md:p-10">
        <Button
          onClick={() => navigate({ to: "/plants" })}
          className='hover:cursor-pointer grow'>
          Plants
          <SproutIcon className="ml-2" />
        </Button>

        <Button
          onClick={() => navigate({ to: "/seeds" })}
          className='hover:cursor-pointer grow'>
          Seeds
          <BeanIcon className="ml-2" />
        </Button>
      </div>
    </div>
  )
}