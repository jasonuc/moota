import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import PlantsList from "@/components/plants-list";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import { getUserNearbyPlants } from "@/services/api/plants";
import { PlantWithDistanceMFromUser } from "@/types/plant";
import { useEffect, useState } from "react";

export default function HomePage() {
  const { user } = useAuth();
  const { latitude, longitude, withinAllowance } = useGeolocation();
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
      <Header />

      <h1 className="text-3xl font-heading mb-2">My Nearby Plants</h1>

      <PlantsList plants={nearbyPlants} showDistanceM={withinAllowance} />

      <div className="pt-20" />

      <DynamicNavigation />
    </div>
  );
}
