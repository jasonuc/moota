import Header from '@/components/header'
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { Progress } from '@/components/ui/progress';
import { createFileRoute } from '@tanstack/react-router'
import { DropletIcon, HouseIcon, MenuIcon, SkullIcon, Triangle } from 'lucide-react';

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
                <img src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${id}&backgroundColor=${"transparent"}`} />

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

            <div className="grid grid-cols-2 grid-rows-2 gap-10">
                <div className="">
                    <p>woe</p>
                    <Progress className='h-4'
                        value={45} />
                </div>
                <div className="">
                    <p>dread</p>
                    <Progress className='h-4'
                        value={65} />
                </div>
                <div className="">
                    <p>frolic</p>
                    <Progress className='h-4'
                        value={80} />
                </div>
                <div className="">
                    <p>malice</p>
                    <Progress className='h-4'
                        value={15} />
                </div>
            </div>

            <Card className="h-56 relative p-0 w-full mt-10">
                <CardHeader className='z-50 absolute inset-0 p-0'>
                    <CardTitle className='p-1.5 w-fit border-b-2 border-r-2 rounded-none text-xs flex items-center justify-center'>
                        <HouseIcon size={15} className='mr-2' />
                        <span>home</span>
                    </CardTitle>
                </CardHeader>
                <CardContent className="size-full flex items-center justify-center">
                    <Triangle />
                </CardContent>
            </Card>

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
