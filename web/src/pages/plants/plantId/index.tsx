import Header from "@/components/header";
import KillPlantButton from "@/components/kill-plant-button";
import MoreButton from "@/components/more-button";
import PlantMap from "@/components/plant-map";
import PlantTempers from "@/components/plant-tempers";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import { PLANT_INTERACTION_RADIUS } from "@/lib/constants";
import {
  cn,
  formatHp,
  formatPlantDate,
  getDicebearThumbsUrl,
  haversineDistance,
  startSentenceWithUppercase,
} from "@/lib/utils";
import { useWaterPlant } from "@/services/mutations/plants";
import { useGetPlant } from "@/services/queries/plants";
import { AxiosError } from "axios";
import { DropletIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router";
import { toast } from "sonner";

export default function IndividualPlantPage() {
  const params = useParams();
  const navigate = useNavigate();
  const { data: plant, error: useGetPlantErr } = useGetPlant(params.plantId);
  const { user } = useAuth();
  const {
    latitude: userLat,
    longitude: userLon,
    withinAllowance,
  } = useGeolocation();
  const waterPlantMtn = useWaterPlant();
  const [isImageVisible, setIsImageVisible] = useState(true);

  useEffect(() => {
    if (useGetPlantErr) {
      const err = useGetPlantErr as AxiosError<{ error: string }>;
      toast.error(err.response?.data.error);
    }
  }, [useGetPlantErr]);

  if (plant?.ownerID && user?.id) {
    if (plant.ownerID != user.id) {
      navigate(`/plants/${plant.id}/public`);
      return;
    }
  }

  const handleWaterPlant = () => {
    waterPlantMtn
      .mutateAsync({
        plantId: plant!.id,
        latitude: userLat!,
        longitude: userLon!,
      })
      .then(() =>
        toast.success(
          <p>
            Watered <em>{plant?.nickname}</em>
          </p>,
          {
            icon: <DropletIcon />,
            position: "top-center",
          }
        )
      )
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

      <div
        className={cn({
          "grid md:flex md:justify-center md:items-center": isImageVisible,
          "grid items-center justify-center": !isImageVisible,
        })}
      >
        <img
          className="mx-auto md:mx-0"
          width={200}
          height={200}
          draggable={false}
          src={getDicebearThumbsUrl(plant?.id)}
          alt={`Avatar for ${plant?.nickname}`}
          onError={(e) => {
            e.currentTarget.style.display = "none";
            setIsImageVisible(false);
          }}
        />

        <div className="flex flex-col items-center justify-center gap-y-2">
          <h1 className="text-2xl font-heading text-center">
            {plant?.nickname}
          </h1>
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

      <PlantMap
        userCoords={{ Lat: userLat ?? 0, Lon: userLon ?? 0 }}
        plantCoords={{
          Lat: plant?.centre.Lat ?? 0,
          Lon: plant?.centre.Lon ?? 0,
        }}
        showUser={withinAllowance}
      />

      <div className="flex flex-col grow justify-end">
        <div className="grid grid-cols-3 gap-x-5 gap-y-5">
          <KillPlantButton
            id={plant?.id || ""}
            nickname={plant?.nickname || ""}
          />

          <Button
            className="md:min-h-12 col-span-2 md:col-span-1 flex items-center justify-center space-x-1.5"
            onClick={handleWaterPlant}
            disabled={
              haversineDistance(
                { Lat: plant?.centre.Lat ?? 0, Lon: plant?.centre.Lon ?? 0 },
                { Lat: userLat ?? 0, Lon: userLon ?? 0 }
              ) > PLANT_INTERACTION_RADIUS
            }
          >
            Water <DropletIcon />
          </Button>

          <MoreButton {...plant!} />
        </div>
      </div>
    </div>
  );
}
