import Header from "@/components/header";
import KillPlantButton from "@/components/kill-plant-button";
import MoreButton from "@/components/more-button";
import PlantMap from "@/components/plant-map";
import PlantTempers from "@/components/plant-tempers";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { formatHp, startSentenceWithUppercase } from "@/lib/utils";
import { getPlant, waterPlant } from "@/services/api/plants";
import { Plant } from "@/types/plant";
import { useGeolocation } from "@uidotdev/usehooks";
import { AxiosError } from "axios";
import { formatDate, isValid, parseJSON } from "date-fns";
import { DropletIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

export default function IndividualPlantPage() {
  const params = useParams();
  const navigate = useNavigate();
  const [plant, setPlant] = useState<Plant>();
  const { user } = useAuth();
  const { latitude, longitude } = useGeolocation();

  useEffect(() => {
    if (!user?.id) return;

    if (params.plantId) {
      getPlant(params.plantId)
        .then((plant) => {
          if (plant.ownerID != user.id) {
            navigate("/home");
            toast.error("You are not allowed to access that page", {
              description: "That plant belongs to another user",
              descriptionClassName: "!text-white",
            });
            return;
          }
          setPlant(plant);
        })
        .catch((err: AxiosError<string>) => {
          toast.error(err.response?.data);
        });
    }
  }, [params.plantId, navigate, user?.id]);

  const formatPlantDate = (dateString: string | undefined) => {
    if (!dateString) return "Unknown";

    try {
      const date = parseJSON(dateString);
      return isValid(date) ? formatDate(date, "dd/MM/yy") : "Invalid date";
    } catch {
      return "Invalid date";
    }
  };

  const handleWaterPlant = () => {
    waterPlant(plant!.id, latitude!, longitude!)
      .then(() => {
        getPlant(params.plantId!)
          .then(setPlant)
          .then(() => {
            toast.success(
              <p>
                Watered <em>{plant?.nickname}</em>
              </p>,
              {
                icon: <DropletIcon />,
                position: "top-center",
              }
            );
          });
      })
      .catch((error: AxiosError<{ error: string }>) => {
        toast.info(
          startSentenceWithUppercase(error.response?.data.error ?? ""),
          {
            description: "This plant was not watered",
            position: "top-center",
          }
        );
      });
  };

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header />

      <div className="grid md:flex grid-cols-2 md:justify-center md:items-center gap-x-10">
        <img
          className="mx-auto md:mx-0"
          width={200}
          height={200}
          draggable={false}
          src={`https://api.dicebear.com/9.x/thumbs/svg?seed=${plant?.id}&backgroundColor=transparent&shapeRotation=-20`}
          alt={`Avatar for ${plant?.nickname}`}
        />

        <div className="flex flex-col items-center justify-center gap-y-2">
          <h1 className="text-2xl font-heading">{plant?.nickname}</h1>
          <small className="italic">{plant?.botanicalName}</small>
          <div className="font-semibold flex gap-x-1.5">
            <p>lv. {plant?.level ?? 0}</p>
            <span>•</span>
            <p>{formatPlantDate(plant?.timePlanted)}</p>
          </div>
          <div className="flex gap-x-4">
            {plant?.soil?.type && (
              <p>
                {`${plant.soil.type
                  .charAt(0)
                  .toUpperCase()}${plant.soil.type.slice(1)} soil`}
              </p>
            )}
            <p>H: {formatHp(plant?.hp ?? 0)}%</p>
          </div>
        </div>
      </div>

      <PlantTempers {...plant?.tempers} />

      <PlantMap />

      <div className="flex flex-col grow justify-end">
        <div className="grid grid-cols-3 gap-x-5 gap-y-5">
          <KillPlantButton
            id={plant?.id || ""}
            nickname={plant?.nickname || ""}
          />

          <Button
            className="md:min-h-12 col-span-2 md:col-span-1 flex items-center justify-center space-x-1.5"
            onClick={handleWaterPlant}
          >
            Water <DropletIcon />
          </Button>

          <MoreButton {...plant!} />
        </div>
      </div>
    </div>
  );
}
