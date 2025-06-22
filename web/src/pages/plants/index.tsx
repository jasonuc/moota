import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import PlantsList from "@/components/plants-list";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import { getAllUserPlants } from "@/services/api/plants";
import { PlantWithDistanceMFromUser } from "@/types/plant";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function AllUserPlantsPage() {
  const { user } = useAuth();
  const { latitude, longitude, withinAllowance } = useGeolocation();
  const [plants, setPlants] = useState<PlantWithDistanceMFromUser[]>();

  useEffect(() => {
    if (!user?.id || !latitude || !longitude) return;

    getAllUserPlants(user.id, latitude, longitude)
      .then(setPlants)
      .catch(() => {
        toast.error("Error occured on the server", {
          description: `Try again later.`,
          descriptionClassName: "!text-white",
        });
      });
  }, [user, latitude, longitude]);

  return (
    <div className="flex flex-col space-y-5 pb-10 grow">
      <Header />

      <h1 className="text-3xl font-heading mb-2">
        My Plants {plants?.length && `(${plants?.length})`}
      </h1>

      <PlantsList
        plants={plants}
        maxPlants={plants?.length}
        showDistanceM={withinAllowance}
      />
      <DynamicNavigation />
    </div>
  );
}
