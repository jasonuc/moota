import Plant from "./plant";
import { Link } from "@tanstack/react-router";
import { Button } from "./ui/button";

type PlantProps = {
    maxPlants?: number;
    plants: {
        nickname: string;
        botanicalName: string;
        hp: number;
        distance: number;
    }[]
}

// TODO: Complete empty state
export default function PlantsList({ plants, maxPlants = 4 }: PlantProps) {
    return (
        <div className="flex flex-col grow">
            {plants.length ?
                (
                    <div className="flex flex-col space-y-5 md:space-y-7">
                        {plants.slice(0, maxPlants).map((p, key) => (
                            <Plant key={key} {...p} />
                        ))}
                    </div>
                ) :
                (
                    <div className="flex flex-col grow items-center justify-center gap-y-5 py-40">
                        <div className="flex flex-col items-center gap-y-1.5">
                            <h3 className="text-xl font-heading">{"Hello, there! 👋"}</h3>
                            <p>You don't have any plants yet. Try adding some!</p>
                        </div>
                        <div className="flex flex-col md:flex-row gap-x-5 gap-y-5">
                            <Button>
                                Plant a seed
                            </Button>
                            <Link to="/">
                                <Button variant="neutral">
                                    <p className="text-blue-600">Learn more</p>
                                </Button>
                            </Link>
                        </div>
                    </div>
                )}
        </div>
    )
}