import Plant from "./plant";

export default function MyNearbyPlants() {
    return (
        <div className="flex flex-col">
            <h1 className="text-3xl font-heading mb-5">My Nearby Plants</h1>

            <div className="flex flex-col space-y-4 h-[40rem]">
                {[1,2,3,4].map((_, key) => (
                    <Plant key={key} />
                ))}
            </div>
        </div>
    )
}