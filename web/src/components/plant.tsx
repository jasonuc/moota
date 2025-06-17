import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { formatDistance } from "@/lib/utils";
import type { PlantWithDistanceMFromUser } from "@/types/plant";
import { Heart, LocateFixed } from "lucide-react";
import { Link } from "react-router";
import { Button } from "./ui/button";

type PlantProps = PlantWithDistanceMFromUser;

export default function Plant({
  id,
  nickname,
  botanicalName,
  hp,
  distanceM,
}: PlantProps) {
  return (
    <Link to={`/plants/${id}`}>
      <Button asChild className="relative overflow-hidden group h-fit">
        <Card className="gap-y-1.5 bg-background flex w-full">
          <CardContent className="flex justify-end w-full gap-x-2 p-0 pt-1.5">
            <div className="flex items-center gap-x-1.5">
              <Heart size={15} />
              {hp}%
            </div>
            <div className="flex items-center gap-x-1.5">
              <LocateFixed size={15} />
              {formatDistance(distanceM!)}
            </div>
          </CardContent>
          <CardHeader className="w-full p-0 pb-3">
            <CardTitle className="text-lg font-bold">{nickname}</CardTitle>
            <CardDescription className="italic">
              {botanicalName}
            </CardDescription>
          </CardHeader>

          <img
            src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${nickname}&backgroundColor=${"transparent"}&shapeRotation=-20`}
            alt="avatar"
            draggable={false}
            className="size-20 pointer-events-none absolute -bottom-1/4 right-0 group-hover:-bottom-5 group-active:-bottom-5 transition-all duration-100 rounded-l-md"
          />
        </Card>
      </Button>
    </Link>
  );
}
