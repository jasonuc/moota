import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import PlantsList from "@/components/plants-list";
import { useAuth } from "@/hooks/use-auth";
import { useGeolocation } from "@/hooks/use-geolocation";
import { useGetUserPlants } from "@/services/queries/plants";
import { useEffect } from "react";
import { toast } from "sonner";

export default function UserPlantsPage() {
  const { user } = useAuth();
  const { latitude, longitude, withinAllowance } = useGeolocation();
  const {
    data: plants,
    error: useGetUserPlantsErr,
    isPlaceholderData,
  } = useGetUserPlants(user?.id, latitude, longitude);

  useEffect(() => {
    if (useGetUserPlantsErr) {
      toast.error("Error occured on the server", {
        description: `Try again later.`,
        descriptionClassName: "!text-white",
      });
    }
  }, [useGetUserPlantsErr]);

  return (
    <div className="flex flex-col space-y-5 pb-10 grow">
      <Header />

      <h1 className="text-3xl font-heading mb-2">
        My Plants {(plants?.length ?? -1) > 0 && `(${plants?.length})`}
      </h1>

      <PlantsList
        plants={plants}
        maxPlants={plants?.length}
        showDistanceM={withinAllowance}
        muteDistanceM={isPlaceholderData}
      />
      <DynamicNavigation />
    </div>
  );
}
