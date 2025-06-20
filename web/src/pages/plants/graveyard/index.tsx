import DynamicNavigation from "@/components/dynamic-navigation";
import Header from "@/components/header";
import Headstone from "@/components/headstone";
import { useAuth } from "@/hooks/use-auth";
import { getAllUserDeceasedPlants } from "@/services/api/plants";
import { Plant } from "@/types/plant";
import { useEffect, useState } from "react";
import { toast } from "sonner";

export default function PlantGraveyard() {
  const { user } = useAuth();
  const [plants, setPlants] = useState<Plant[]>();

  useEffect(() => {
    if (!user?.id) return;

    getAllUserDeceasedPlants(user.id)
      .then(setPlants)
      .catch(() => {
        toast.error("Error occured on the server", {
          description: `Try again later.`,
          descriptionClassName: "!text-white",
        });
      });
  }, [user]);

  return (
    <div className="flex flex-col space-y-5 pb-10 grow">
      <Header />

      <h1 className="text-3xl font-heading mb-2">
        Graveyard ({plants?.length || 0})
      </h1>

      {plants?.length ? (
        <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
          {plants.map((plant) => (
            <Headstone key={plant.id} plant={plant} />
          ))}
        </div>
      ) : (
        <div className="flex flex-col grow items-center justify-center gap-y-5 py-40">
          <div className="flex flex-col items-center gap-y-1.5">
            <h3 className="text-xl font-heading">{"No plants here yet 🌱"}</h3>
            <p className="text-center text-slate-600">
              Your plants are still alive and well!
              <br />
              Keep taking care of them to avoid visits here.
            </p>
          </div>
        </div>
      )}

      <DynamicNavigation />
    </div>
  );
}
