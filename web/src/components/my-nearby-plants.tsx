import Plant from "./plant";
import { Link } from "@tanstack/react-router";
import { Button } from "./ui/button";

// TODO: Complete empty state
export default function MyNearbyPlants() {
    const nearbyPlants = [
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
    ]

    return (
        <div className="flex flex-col grow">
            <h1 className="text-3xl font-heading mb-5">My Nearby Plants</h1>

            {nearbyPlants.length ?
                (
                    <div className="flex flex-col space-y-5 md:space-y-7">
                        {nearbyPlants.map((p, key) => (
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