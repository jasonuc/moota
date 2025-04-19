import Plant from "./plant";

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
        <div className="flex flex-col">
            <h1 className="text-3xl font-heading mb-5">My Nearby Plants</h1>

            <div className="flex flex-col space-y-5 md:space-y-7">
                {nearbyPlants.map((p, key) => (
                    <Plant key={key} {...p} />
                ))}
            </div>
        </div>
    )
}