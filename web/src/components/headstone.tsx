import type { Plant } from "@/types/plant";
import {
  Card,
  CardHeader,
  CardTitle,
  CardDescription,
  CardContent,
} from "./ui/card";
import { formatCoordinates, formatRelativeTime } from "@/lib/utils";
import { GhostIcon } from "lucide-react";

function generateDeceasedPlantEpitaph(
  botanicalName: string,
  nickname: string,
): string {
  const templates = [
    "{nickname} was the most magnificent {botanical} that ever graced this earth.",
    "{nickname} stood as the legendary {botanical} of ancient tales.",
    "{nickname} was the fierce {botanical} that conquered hearts across the land.",
    "{nickname} became the mystical {botanical} whispered about in garden lore.",
    "{nickname} was the rebellious {botanical} that refused to follow the rules.",
    "{nickname} emerged as the wise {botanical} that knew all nature's secrets.",
    "{nickname} was the adventurous {botanical} that explored every corner of the world.",
    "{nickname} stood proud as the noble {botanical} of royal gardens.",
    "{nickname} was the enchanted {botanical} that brought magic to every season.",
    "{nickname} became the legendary {botanical} that inspired countless stories.",
    "{nickname} was the bravest {botanical} that ever faced the storms of life.",
    "{nickname} became the mysterious {botanical} that haunted gardeners' dreams.",
    "{nickname} was the immortal {botanical} that transcended all seasons.",
    "{nickname} emerged as the cunning {botanical} that outsmarted nature itself.",
    "{nickname} was the glorious {botanical} that ruled the botanical kingdom.",
    "{nickname} became the tragic {botanical} whose beauty was too pure for this world.",
    "{nickname} was the defiant {botanical} that bloomed against all odds.",
    "{nickname} stood as the eternal {botanical} in the halls of plant legends.",
    "{nickname} was the charismatic {botanical} that charmed every creature nearby.",
    "{nickname} became the sacred {botanical} worshipped by ancient civilizations.",
  ];

  const hash = (botanicalName + nickname)
    .split("")
    .reduce((acc, char) => acc + char.charCodeAt(0), 0);

  const templateIndex = hash % templates.length;
  const selectedTemplate = templates[templateIndex];

  return selectedTemplate
    .replace("{nickname}", nickname)
    .replace("{botanical}", botanicalName);
}

type HeadstoneProps = {
  plant: Plant;
};

export default function Headstone({
  plant: { nickname, timePlanted, centre, timeOfDeath, botanicalName },
}: HeadstoneProps) {
  return (
    <Card className="relative">
      <CardHeader>
        <GhostIcon className="mb-2 absolute right-2 top-2" />
        <CardTitle>{nickname}</CardTitle>
        <CardDescription>
          Time of Death: {formatRelativeTime(timeOfDeath)}
        </CardDescription>
      </CardHeader>
      <CardContent>
        Their journey began {formatRelativeTime(timePlanted)} {" at "}
        <span>{formatCoordinates(centre.Lat, centre.Lon)}</span>.{" "}
        {generateDeceasedPlantEpitaph(botanicalName, nickname)}
      </CardContent>
    </Card>
  );
}
