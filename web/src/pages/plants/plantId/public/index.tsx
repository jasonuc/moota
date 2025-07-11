import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import { Button } from "@/components/ui/button";
import { Progress } from "@/components/ui/progress";
import { SERVER_404_MESSAGE } from "@/lib/constants";
import {
  cn,
  formatHp,
  formatPlantDate,
  getDicebearThumbsUrl,
  startSentenceWithUppercase,
} from "@/lib/utils";
import { useGetPlant } from "@/services/queries/plants";
import { useGetUsernameFromUserId } from "@/services/queries/user";
import { AxiosError } from "axios";
import { ShareIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Link, useNavigate, useParams } from "react-router";
import { toast } from "sonner";

export default function PublicPlantPage() {
  const params = useParams();
  const { data: plant, error: useGetPlantErr } = useGetPlant(params.plantId);
  const { data: ownerUsername } = useGetUsernameFromUserId(plant?.ownerID);
  const [isImageVisible, setIsImageVisible] = useState(true);
  const navigate = useNavigate();

  useEffect(() => {
    if (useGetPlantErr) {
      const err = useGetPlantErr as AxiosError<{ error: string }>;
      if (err.response?.data.error === SERVER_404_MESSAGE) {
        navigate("/home");
        toast.error("Plant does not exist", {
          description: err.response?.data.error,
          descriptionClassName: "!text-white",
        });
        return;
      }
      toast.error("A problem has occured on the server");
    }
  }, [useGetPlantErr, navigate]);

  if (useGetPlantErr) return null;

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header />

      <div
        className={cn("w-full md:max-w-md md:mx-auto", {
          "h-[80%] md:h-full flex items-center justify-center": !isImageVisible,
        })}
      >
        <div className="bg-white/20 backdrop-blur-3xl border-2 border-white/80 rounded-lg shadow-lg px-6 py-10">
          <img
            className="mx-auto mb-6"
            width={200}
            height={200}
            draggable={false}
            src={getDicebearThumbsUrl(plant?.id)}
            onError={(e) => {
              e.currentTarget.style.display = "none";
              setIsImageVisible(false);
            }}
            alt={`Avatar for ${plant?.nickname}`}
          />

          <div className="text-center space-y-2 mb-6">
            <h1 className="text-3xl font-heading">{plant?.nickname}</h1>
            <p className="italic text-sm">{plant?.botanicalName}</p>

            <p className="text-sm mt-3">
              Owned by:{" "}
              <Link to={`/profile/${ownerUsername}`} className="underline">
                {ownerUsername}
              </Link>
            </p>
          </div>

          <div className="grid grid-cols-2 gap-4 mb-6">
            <div className="text-center">
              <p className="text-sm font-semibold">Level</p>
              <p className="text-2xl font-bold">{plant?.level}</p>
            </div>
            <div className="text-center">
              <p className="text-sm font-semibold">XP</p>
              <p className="text-2xl font-bold">{plant?.xp}</p>
            </div>
            <div className="text-center">
              <p className="text-sm font-semibold">Health</p>
              <p className="text-2xl font-bold">{formatHp(plant?.hp ?? 0)}%</p>
            </div>
            <div className="text-center">
              <p className="text-sm font-semibold">Planted</p>
              <p className="text-2xl">{formatPlantDate(plant?.timePlanted)}</p>
            </div>
          </div>

          {plant?.soil?.type && (
            <div className="text-center mb-6">
              <p className="">
                {startSentenceWithUppercase(plant.soil.type)} Soil
              </p>
            </div>
          )}

          <div>
            <h3 className="font-heading text-center mb-4">Tempers</h3>
            <div className="grid grid-cols-2 gap-x-10 gap-y-4">
              {[
                { label: "woe", value: plant?.tempers?.woe },
                { label: "dread", value: plant?.tempers?.dread },
                { label: "frolic", value: plant?.tempers?.frolic },
                { label: "malice", value: plant?.tempers?.malice },
              ].map(({ label, value }) => (
                <div key={label}>
                  <p className="text-sm italic capitalize mb-2">{label}</p>
                  <Progress
                    className="h-5 rounded-md md:rounded-base"
                    value={typeof value === "number" ? (value / 5) * 100 : 0}
                  />
                  <p className="text-xs text-right mt-1">
                    {typeof value === "number" ? value : 0}/5
                  </p>
                </div>
              ))}
            </div>
          </div>
        </div>
      </div>

      <div className="fixed left-0 bottom-0 w-full p-5 flex justify-start">
        <Button
          className="rounded-full size-12 z-50"
          onClick={() => {
            navigator.clipboard
              .writeText(window.location.href)
              .then(() => toast.success("Copied link to plant card!"));
          }}
        >
          <ShareIcon />
        </Button>
        <DynamicNavigation />
      </div>
    </div>
  );
}
