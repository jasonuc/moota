import Header from "@/components/header";
import PlantsList from "@/components/plants-list";
import { Button } from "@/components/ui/button";
import { useAuth } from "@/hooks/use-auth";
import { getUserNearbyPlants } from "@/services/api/plants";
import { PlantWithDistanceMFromUser } from "@/types/plant";
import { useGeolocation } from "@uidotdev/usehooks";
import { BeanIcon, SproutIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

export default function HomePage() {
  const { user } = useAuth();
  const navigate = useNavigate();
  const { latitude, longitude } = useGeolocation();
  const [nearbyPlants, setNearbyPlants] = useState<
    PlantWithDistanceMFromUser[] | undefined
  >();

  useEffect(() => {
    if (!user?.id || !latitude || !longitude) {
      return;
    }

    getUserNearbyPlants(user.id, latitude, longitude)
      .then(setNearbyPlants)
      .catch(console.error);
  }, [user?.id, latitude, longitude]);

  return (
    <div className="flex flex-col space-y-5 grow">
      <Header seedCount={10} />

      <h1 className="text-3xl font-heading mb-5">My Nearby Plants</h1>

      <PlantsList plants={nearbyPlants} />

      <div className="pt-20" />

      <div className="fixed left-0 bottom-0 w-full flex gap-x-5 p-5 md:p-10">
        <Button
          onClick={() => navigate({ pathname: "/plants" })}
          className="hover:cursor-pointer grow md:min-h-12"
        >
          Plants
          <SproutIcon className="ml-2" />
        </Button>

        <Button
          onClick={() => navigate({ pathname: "/seeds" })}
          className="hover:cursor-pointer grow md:min-h-12"
        >
          Seeds
          <BeanIcon className="ml-2" />
        </Button>
      </div>
    </div>
  );
}
