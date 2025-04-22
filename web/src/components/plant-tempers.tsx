import { Progress } from "./ui/progress";

type PlantTempersProps = {
    woe: number;
    frolic: number;
    malice: number;
    dread: number;
}

// TODO: Styling not so perfect when in between sm and md
export default function PlantTempers({ woe, dread, frolic, malice }: PlantTempersProps) {
    return (
        <div>
            <h3 className='font-heading mb-3'>Tempers</h3>
            <div className="grid grid-cols-2 grid-rows-2 gap-y-5 gap-x-10 md:gap-10">
                <div className="">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-3 text-left'>woe</p>
                    <Progress className='h-5 rounded-md md:rounded-base rounde md:min-h-8'
                        value={(woe / 5) * 100} />
                </div>
                <div className="">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-3 text-left'>dread</p>
                    <Progress className='h-5 rounded-md md:rounded-base rounde md:min-h-8'
                        value={(dread / 5) * 100} />
                </div>
                <div className="">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-3 text-left'>frolic</p>
                    <Progress className='h-5 rounded-md md:rounded-base rounde md:min-h-8'
                        value={(frolic / 5) * 100} />
                </div>
                <div className="">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-3 text-left'>malice</p>
                    <Progress className='h-5 rounded-md md:rounded-base rounde md:min-h-8'
                        value={(malice / 5) * 100} />
                </div>
            </div>
        </div>
    )
}