import Header from "@/components/header";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { startSentenceWithUppercase } from "@/lib/utils";
import { getUserSeeds, plantSeed } from "@/services/api/seeds";
import { Seed, SeedGroup } from "@/types/seed";
import { useGeolocation } from "@uidotdev/usehooks";
import { AxiosError } from "axios";
import { AudioLinesIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

export default function SeedsPage() {
  const [seeds, setSeeds] = useState<SeedGroup[]>();
  const { user } = useAuth();
  const { latitude, longitude } = useGeolocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!user?.id) {
      return;
    }

    getUserSeeds(user.id)
      .then(setSeeds)
      .catch((error) => console.error(error));
  }, [user?.id]);

  const decideSeedToPlant = async (count: number, seeds: Seed[]) => {
    const seedToPlant = seeds[Math.floor(Math.random() * count)];

    plantSeed(seedToPlant.id, latitude!, longitude!)
      .then(() =>
        getUserSeeds(user!.id).then(() => navigate("/plants/unactivated"))
      )
      .catch((error: AxiosError<{ error: string }>) => {
        console.log(error.response?.data.error);
        toast.error(
          startSentenceWithUppercase(error.response?.data.error ?? ""),
          {
            description: "There is another plant in the area",
            descriptionClassName: "!text-white",
          }
        );
      });
  };

  return (
    <div className="flex flex-col space-y-5 pb-10 w-full">
      <Header />

      <h1 className="text-3xl font-heading mb-5">My Seeds</h1>

      <div className="grid grid-cols-3 md:grid-cols-4 gap-5">
        {seeds?.map(({ botanicalName, count, seeds }) => (
          <Button
            asChild
            className="relative h-36 group"
            key={botanicalName}
            onClick={() => decideSeedToPlant(count, seeds)}
          >
            <div className="size-full relative">
              <AudioLinesIcon className="absolute group-active:scale-75 transition-all duration-300 ease-in-out bottom-0 left-0 rotate-45" />
              <AudioLinesIcon className="absolute group-active:scale-75 transition-all duration-300 ease-in-out bottom-0 right-0 -rotate-45" />
              <AudioLinesIcon className="absolute group-active:scale-75 transition-all duration-300 ease-in-out top-0 right-0 rotate-45" />
              {!(count > 1) && (
                <AudioLinesIcon className="absolute group-active:scale-75 transition-all duration-300 ease-in-out top-0 left-0 -rotate-45" />
              )}

              {count > 1 && (
                <small className="absolute left-1 -top-2 bg-background px-2 rounded-full">
                  x{count}
                </small>
              )}
              <p className="italic text-wrap text-center">{botanicalName}</p>
            </div>
          </Button>
        ))}
      </div>
    </div>
  );
}
