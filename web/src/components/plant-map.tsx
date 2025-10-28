import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { cn, formatDistance, haversineDistance } from "@/lib/utils";
import { Coordinates } from "@/types/coordinates";
import { Flower2Icon, LocateFixedIcon, PersonStandingIcon } from "lucide-react";
import {
  AlertDialog,
  AlertDialogAction,
  AlertDialogCancel,
  AlertDialogContent,
  AlertDialogDescription,
  AlertDialogFooter,
  AlertDialogHeader,
  AlertDialogTitle,
  AlertDialogTrigger,
} from "./ui/alert-dialog";

type PlantMapProps = {
  plantCoords: Coordinates;
  userCoords?: Coordinates;
  showUser?: boolean;
  plantRadiusM?: number;
};

export default function PlantMap({
  plantCoords,
  userCoords,
  showUser = false,
  plantRadiusM = 0,
}: PlantMapProps) {
  const mapWidth = 400;
  const mapHeight = 200;

  const plantX = ((plantCoords.Lon + 180) / 360) * mapWidth;
  const plantY = ((90 - plantCoords.Lat) / 180) * mapHeight;

  const userX = userCoords ? ((userCoords.Lon + 180) / 360) * mapWidth : 0;
  const userY = userCoords ? ((90 - userCoords.Lat) / 180) * mapHeight : 0;

  const distance = userCoords ? haversineDistance(plantCoords, userCoords) : 0;
  const isAtPlant = distance <= plantRadiusM;

  return (
    <AlertDialog>
      <AlertDialogTrigger>
        <Card className="grow max-h-80 h-56 min-h-56 relative p-0 w-full mt-10 md:mt-5 overflow-hidden">
          <CardHeader className="z-50 absolute inset-0 p-0">
            <CardTitle
              className={cn(
                "p-1 w-fit border-b-2 border-r-2 rounded-none rounded-br-md text-xs flex items-center justify-center font-semibold bg-white/90 backdrop-blur-sm",
                { isAtPlant: "text-green-600" },
              )}
            >
              <LocateFixedIcon size={15} className="mr-1" />
              {isAtPlant && showUser ? (
                <span className=" ">You are here</span>
              ) : (
                <span>Away from plant</span>
              )}
            </CardTitle>
          </CardHeader>
          <CardContent className="size-full relative p-0">
            <div className="absolute inset-0 bg-gradient-to-br from-green-100 to-blue-200">
              <div className="absolute inset-0 opacity-20">
                {[...Array(8)].map((_, i) => (
                  <div
                    key={`h-${i}`}
                    className="absolute w-full h-px bg-gray-400"
                    style={{ top: `${i * 12.5}%` }}
                  />
                ))}
                {[...Array(12)].map((_, i) => (
                  <div
                    key={`v-${i}`}
                    className="absolute h-full w-px bg-gray-400"
                    style={{ left: `${i * 8.33}%` }}
                  />
                ))}
              </div>

              {showUser && userCoords && !isAtPlant && (
                <>
                  <svg className="absolute inset-0 w-full h-full">
                    <line
                      x1={`${Math.max(
                        5,
                        Math.min(95, (plantX / mapWidth) * 100),
                      )}%`}
                      y1={`${Math.max(
                        5,
                        Math.min(95, (plantY / mapHeight) * 100),
                      )}%`}
                      x2={`${Math.max(
                        5,
                        Math.min(95, (userX / mapWidth) * 100),
                      )}%`}
                      y2={`${Math.max(
                        5,
                        Math.min(95, (userY / mapHeight) * 100),
                      )}%`}
                      stroke="#3b82f6"
                      strokeWidth="2"
                      strokeDasharray="4 4"
                    />
                  </svg>

                  <div
                    className="absolute bg-blue-500 text-white text-xs px-2 py-1 rounded-full font-semibold"
                    style={{
                      left: `${
                        (Math.max(5, Math.min(95, (plantX / mapWidth) * 100)) +
                          Math.max(5, Math.min(95, (userX / mapWidth) * 100))) /
                        2
                      }%`,
                      top: `${
                        (Math.max(5, Math.min(95, (plantY / mapHeight) * 100)) +
                          Math.max(
                            5,
                            Math.min(95, (userY / mapHeight) * 100),
                          )) /
                          2 -
                        5
                      }%`,
                      transform: "translate(-50%, -50%)",
                    }}
                  >
                    {formatDistance(distance)}
                  </div>

                  <div
                    className="absolute transform -translate-x-1/2 -translate-y-full z-10"
                    style={{
                      left: `${Math.max(
                        5,
                        Math.min(95, (userX / mapWidth) * 100),
                      )}%`,
                      top: `${Math.max(
                        5,
                        Math.min(95, (userY / mapHeight) * 100),
                      )}%`,
                    }}
                  >
                    <PersonStandingIcon
                      className="text-indigo-500 drop-shadow-lg"
                      size={20}
                    />
                  </div>
                </>
              )}

              <div
                className="absolute transform -translate-x-1/2 -translate-y-full z-10"
                style={{
                  left: `${Math.max(
                    5,
                    Math.min(95, (plantX / mapWidth) * 100),
                  )}%`,
                  top: `${Math.max(
                    5,
                    Math.min(95, (plantY / mapHeight) * 100),
                  )}%`,
                }}
              >
                <Flower2Icon
                  className="text-green-500 drop-shadow-lg"
                  size={24}
                />
              </div>

              <div className="absolute bottom-2 right-2 bg-black/70 text-white text-xs px-2 py-1 rounded font-mono">
                {plantCoords.Lat.toFixed(4)}, {plantCoords.Lon.toFixed(4)}
              </div>
            </div>
          </CardContent>
        </Card>
      </AlertDialogTrigger>
      <AlertDialogContent>
        <AlertDialogHeader>
          <AlertDialogTitle>Open in Google Maps</AlertDialogTitle>
          <AlertDialogDescription>
            Get directions to your plant location. This will open in a new tab
            or the Google Maps app.
          </AlertDialogDescription>
        </AlertDialogHeader>
        <AlertDialogFooter>
          <AlertDialogCancel>Cancel</AlertDialogCancel>
          <AlertDialogAction asChild>
            <a
              href={`https://www.google.com/maps?q=${plantCoords.Lat},${plantCoords.Lon}`}
              target="_blank"
              rel="noopener noreferrer"
            >
              Open Google Maps
            </a>
          </AlertDialogAction>
        </AlertDialogFooter>
      </AlertDialogContent>
    </AlertDialog>
  );
}
