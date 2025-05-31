import Header from "@/components/header";
import PlantMap from "@/components/plant-map";
import PlantTempers from "@/components/plant-tempers";
import { Button } from "@/components/ui/button";
import { createFileRoute } from "@tanstack/react-router";
import { DropletIcon, MenuIcon, SkullIcon } from "lucide-react";

export const Route = createFileRoute("/plants/$plantId/")({
  component: RouteComponent,
});

// TODO: This component should run an auth check so other non-owners do not have access to this page or at least not in it's entirety
function RouteComponent() {
  const { id, nickname, botanicalName, level, plantedAt, soilType, hp } = {
    id: "1",
    nickname: "Sproutlet",
    botanicalName: "Monstera deliciosa",
    level: 3,
    soilType: "Loam",
    hp: 78,
    plantedAt: "12/02/24",
  };

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header seedCount={10} />

      <div className="grid md:flex grid-cols-2 md:justify-center md:items-center gap-x-10">
        <img
          className="mx-auto md:mx-0"
          width={200}
          height={200}
          draggable={false}
          src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${id}&backgroundColor=${"transparent"}&shapeRotation=-20`}
        />

        <div className="flex flex-col items-center justify-center gap-y-2">
          <h1 className="text-2xl font-heading">{nickname}</h1>
          <small className="italic">{botanicalName}</small>
          <div className="font-semibold flex gap-x-1.5">
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

      <div className="flex flex-col grow justify-end">
        <div className="grid grid-cols-3 gap-x-5 gap-y-5">
          <Button className="md:min-h-12 col-span-1 flex items-center justify-center space-x-1.5 bg-red-400">
            Kill <SkullIcon />
          </Button>

          <Button className="md:min-h-12 col-span-2 md:col-span-1 flex items-center justify-center space-x-1.5">
            Water <DropletIcon />
          </Button>

          <Button
            className="md:min-h-12 col-span-3 md:col-span-1 flex items-center justify-center space-x-1.5"
            variant="neutral"
          >
            more <MenuIcon />
          </Button>
        </div>
      </div>
    </div>
  );
}
