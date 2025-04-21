import Header from '@/components/header'
import PlantMap from '@/components/plant-map';
import PlantTempers from '@/components/plant-tempers';
import { Button } from '@/components/ui/button';
import { createFileRoute } from '@tanstack/react-router'
import { DropletIcon, MenuIcon, SkullIcon } from 'lucide-react';

export const Route = createFileRoute('/plants/$plantId/')({
    component: RouteComponent,
})

// TODO: This component should run an auth check so other non-owners do not have access to this page
function RouteComponent() {
    const { id, nickname, botanicalName, level, plantedAt, soilType, hp } = {
        id: "1",
        nickname: "Sproutlet",
        botanicalName: "Monstera deliciosa",
        level: 3,
        soilType: "Loam",
        hp: 78,
        plantedAt: "12/02/24"
    };

    return (
        <div className='flex flex-col space-y-5 grow'>
            <Header seedCount={10} />

            <div className='grid grid-cols-2'>
                <img className='mx-auto' width={200} height={200} src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${id}&backgroundColor=${"transparent"}`} />

                <div className="flex flex-col items-center justify-center gap-y-2">
                    <h1 className="text-2xl font-heading">{nickname}</h1>
                    <small className='italic'>{botanicalName}</small>
                    <div className='font-semibold flex gap-x-1.5'>
                        <p>lv. {level}</p>
                        <span>•</span>
                        <p>{plantedAt}</p>
                    </div>
                    <div className="flex gap-x-4">
                        <p>{soilType}</p>
                        <p>H: {hp}%</p>
                    </div>
                </div>
            </div>

            <PlantTempers woe={4} frolic={3} dread={1} malice={2} />

            <PlantMap />

            <div className='flex flex-col gap-y-5 grow justify-end'>
                <div className='grid grid-cols-3 gap-x-5'>
                    <Button className='col-span-1 flex items-center justify-center space-x-1.5 bg-red-400'>
                        Kill <SkullIcon />
                    </Button>

                    <Button className='col-span-2 flex items-center justify-center space-x-1.5'>
                        Water <DropletIcon />
                    </Button>
                </div>
                <Button className='flex items-center justify-center space-x-1.5' variant='neutral'>
                    more <MenuIcon />
                </Button>
            </div>
        </div>
    )
}
