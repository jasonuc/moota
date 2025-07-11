import { PlantWithDistanceMFromUser } from "@/types/plant";
import { Link } from "react-router";
import LearnMoreButton from "./learn-more-button";
import PlantListItem from "./plant-list-item";
import { Button } from "./ui/button";

type PlantsListProps = {
  maxPlants?: number;
  showDistanceM: boolean;
  muteDistanceM: boolean;
  plants: PlantWithDistanceMFromUser[] | undefined;
};

export default function PlantsList({
  plants = [],
  showDistanceM,
  muteDistanceM = false,
  maxPlants = 4,
}: PlantsListProps) {
  return (
    <div className="flex flex-col grow">
      {plants.length ? (
        <div className="flex flex-col space-y-5 md:space-y-7">
          {plants.slice(0, maxPlants).map((p, key) => (
            <PlantListItem
              key={key}
              {...p}
              showDistanceM={showDistanceM}
              muteDistanceM={muteDistanceM}
            />
          ))}
        </div>
      ) : (
        <div className="flex flex-col grow items-center justify-center gap-y-5 py-40">
          <div className="flex flex-col items-center gap-y-1.5">
            <h3 className="text-xl font-heading">{"Hello, there! 👋"}</h3>
            <p>You don't have any plants yet. Try planting some!</p>
          </div>
          <div className="flex flex-col md:flex-row gap-x-5 gap-y-5">
            <Link to="/seeds">
              <Button>Plant a seed</Button>
            </Link>
            <LearnMoreButton />
          </div>
        </div>
      )}

      {plants.length < 4 && plants.length != 0 && (
        <div className="flex flex-col items-center gap-y-1.5 py-10">
          <h3 className="text-lg font-heading">{"Not enough plants :("}</h3>
          <Link
            to="/seeds"
            className="text-blue-700 underline underline-offset-2"
          >
            <p>Try and make it to 4 plants. It's more fun that way!</p>
          </Link>
        </div>
      )}
    </div>
  );
}
