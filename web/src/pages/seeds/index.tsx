import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import NoSeeds from "@/components/no-seeds";
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
} from "@/components/ui/alert-dialog";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import { startSentenceWithUppercase } from "@/lib/utils";
import { usePlantSeed } from "@/services/mutations/seeds";
import { useGetUserSeeds } from "@/services/queries/seeds";
import { Seed } from "@/types/seed";
import { AxiosError } from "axios";
import { AudioLinesIcon, SproutIcon, XIcon } from "lucide-react";
import { useEffect } from "react";
import { useNavigate } from "react-router";
import { toast } from "sonner";

export default function SeedsPage() {
  const { user } = useAuth();
  const { latitude, longitude } = useGeolocation();
  const navigate = useNavigate();

  const { withinAllowance } = useGeolocation();
  const { data: seeds, error: useGetUserSeedsErr } = useGetUserSeeds(user?.id);
  const plantSeedMtn = usePlantSeed();

  useEffect(() => {
    if (useGetUserSeedsErr) {
      const err = useGetUserSeedsErr as AxiosError<{ error: string }>;
      toast.error("Error occured on the server", {
        description: `Seeds could not be fetched. ${startSentenceWithUppercase(
          err.response?.data.error ?? "",
        )}`,
        descriptionClassName: "!text-white",
      });
    }
  });

  const decideSeedToPlant = async (count: number, seeds: Seed[]) => {
    const seedToPlant = seeds[Math.floor(Math.random() * count)];

    plantSeedMtn
      .mutateAsync({
        seedId: seedToPlant.id,
        lat: latitude!,
        lon: longitude!,
      })
      .then(() => navigate("/plants"))
      .catch((error: AxiosError<{ error: string }>) => {
        toast.error(
          startSentenceWithUppercase(error.response?.data.error ?? ""),
          {
            description: "There is another plant in the area",
            descriptionClassName: "!text-white",
            position: "top-center",
          },
        );
      });
  };

  return (
    <div className="flex flex-col space-y-5 pb-10 w-full">
      <Header />

      <h1 className="text-3xl font-heading mb-5">
        My Seeds {seeds?.length && `(${seeds.length})`}
      </h1>

      {seeds?.length && (
        <div className="grid grid-cols-3 md:grid-cols-4 gap-5">
          {seeds?.map(({ botanicalName, count, seeds }) => (
            <AlertDialog key={botanicalName}>
              <AlertDialogTrigger asChild>
                <Button asChild className="relative h-36 group">
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
                    <p className="italic text-wrap text-center">
                      {botanicalName}
                    </p>
                  </div>
                </Button>
              </AlertDialogTrigger>
              <AlertDialogContent>
                <AlertDialogHeader>
                  <AlertDialogTitle>
                    {withinAllowance
                      ? "Plant your seed here?"
                      : "Device location too imprecise"}
                  </AlertDialogTitle>
                  <AlertDialogDescription>
                    {withinAllowance
                      ? "By planting here, you agree to visit regularly to care for your plant."
                      : "For better accuracy, try using a mobile device. If you're already on mobile, move outside or closer to a window for better GPS signal."}
                  </AlertDialogDescription>
                </AlertDialogHeader>
                <AlertDialogFooter>
                  <AlertDialogCancel className="hover:cursor-pointer md:min-h-12">
                    <XIcon />
                    {withinAllowance ? "No" : "Close"}
                  </AlertDialogCancel>
                  {withinAllowance && (
                    <AlertDialogAction
                      onClick={() => decideSeedToPlant(count, seeds)}
                      className="hover:cursor-pointer md:min-h-12 col-span-1 flex items-center justify-center space-x-1.5"
                    >
                      <SproutIcon />
                      Yes
                    </AlertDialogAction>
                  )}
                </AlertDialogFooter>
              </AlertDialogContent>
            </AlertDialog>
          ))}
        </div>
      )}

      {!seeds?.length && <NoSeeds />}
      <DynamicNavigation />
    </div>
  );
}
