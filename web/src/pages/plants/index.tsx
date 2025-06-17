import Header from "@/components/header";
import PlantsList from "@/components/plants-list";
import { useAuth } from "@/hooks/use-auth";
import { getAllUserPlants } from "@/services/api/plants";
import { PlantWithDistanceMFromUser } from "@/types/plant";
import { useGeolocation } from "@uidotdev/usehooks";
import { useEffect, useState } from "react";

export default function AllUserPlantsPage() {
  const { user } = useAuth();
  const { latitude, longitude } = useGeolocation();
  const [plants, setPlants] = useState<
    PlantWithDistanceMFromUser[] | undefined
  >();

  useEffect(() => {
    if (!user?.id || !latitude || !longitude) {
      return;
    }

    getAllUserPlants(user.id, latitude, longitude)
      .then(setPlants)
      .catch(console.error);
  }, [user, latitude, longitude]);

  return (
    <div className="flex flex-col space-y-5 pb-10 grow">
      <Header />

      <h1 className="text-3xl font-heading mb-5">
        My Plants ({plants?.length})
      </h1>

      <PlantsList plants={plants} maxPlants={plants?.length} />
    </div>
  );
}
