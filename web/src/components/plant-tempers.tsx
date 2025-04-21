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
            <div className="grid grid-cols-4 md:grid-cols-2 md:grid-rows-2 md:gap-10">
                <div className="even:mt-3 md:even:mt-0">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-10 md:pb-3 text-center md:text-left z-20'>woe</p>
                    <Progress className='h-4 -rotate-90 md:rotate-0 md:min-h-8'
                        value={(woe / 5) * 100} />
                </div>
                <div className="even:mt-3 md:even:mt-0">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-10 md:pb-3 text-center md:text-left z-20'>dread</p>
                    <Progress className='h-4 -rotate-90 md:rotate-0 md:min-h-8'
                        value={(dread / 5) * 100} />
                </div>
                <div className="even:mt-3 md:even:mt-0">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-10 md:pb-3 text-center md:text-left z-20'>frolic</p>
                    <Progress className='h-4 -rotate-90 md:rotate-0 md:min-h-8'
                        value={(frolic / 5) * 100} />
                </div>
                <div className="even:mt-3 md:even:mt-0">
                    <p className='mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-10 md:pb-3 text-center md:text-left z-20'>malice</p>
                    <Progress className='h-4 -rotate-90 md:rotate-0 md:min-h-8'
                        value={(malice / 5) * 100} />
                </div>
            </div>
        </div>
    )
}