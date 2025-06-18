import Header from "@/components/header";
import UnactivatedPlantsList from "@/components/unactivated-plant-list";
import { useAuth } from "@/hooks/use-auth";
import { getAllUserUnactivatedPlants } from "@/services/api/plants";
import { PlantWithDistanceMFromUser } from "@/types/plant";
import { useGeolocation } from "@uidotdev/usehooks";
import { useEffect, useState } from "react";
import { useNavigate } from "react-router";

export default function UnactivatedPlantsPage() {
  const { user } = useAuth();
  const [unactivatedPlants, setUnactivatedPlants] =
    useState<PlantWithDistanceMFromUser[]>();
  const { latitude, longitude } = useGeolocation();
  const navigate = useNavigate();

  useEffect(() => {
    if (!user?.id || !latitude || !longitude) {
      return;
    }

    getAllUserUnactivatedPlants(user.id, latitude, longitude).then(
      (unactivatedPlantCount) => {
        setUnactivatedPlants(unactivatedPlantCount);
        if (unactivatedPlantCount.length <= 0) {
          navigate("/plants");
        }
      }
    );
  }, [user, latitude, longitude, navigate]);

  return (
    <div className="flex flex-col space-y-5 pb-10 grow">
      <Header />

      <div>
        <h1 className="text-3xl font-heading mb-1.5">
          My Unactivated Plants ({unactivatedPlants?.length})
        </h1>
        <p>Activate them before it's too late!</p>
      </div>

      <UnactivatedPlantsList plants={unactivatedPlants ?? []} />
    </div>
  );
}
