import { PlantWithDistanceMFromUser } from "@/types/plant";
import { Link } from "react-router";
import { Button } from "./ui/button";
import UnactivatedPlant from "./unactivated-plant";

interface UnactivatedPlantsListProps {
  plants: PlantWithDistanceMFromUser[];
}

export default function UnactivatedPlantsList({
  plants,
}: UnactivatedPlantsListProps) {
  return (
    <div className="flex flex-col grow">
      {plants.length ? (
        <div className="flex flex-col space-y-5 md:space-y-7">
          {plants.map((p, key) => (
            <UnactivatedPlant key={key} {...p} />
          ))}
        </div>
      ) : (
        <div className="flex flex-col grow items-center justify-center gap-y-5 py-40">
          <div className="flex flex-col items-center gap-y-1.5">
            <h3 className="text-xl font-heading">{"Hello, there! 👋"}</h3>
            <p>You don't have any unactivated plants yet. Try planting some!</p>
          </div>
          <div className="flex flex-col md:flex-row gap-x-5 gap-y-5">
            <Link to="/seeds">
              <Button>Plant a seed</Button>
            </Link>
            <Link to="/">
              <Button variant="neutral">
                <p className="text-blue-600">Learn more</p>
              </Button>
            </Link>
          </div>
        </div>
      )}
    </div>
  );
}
