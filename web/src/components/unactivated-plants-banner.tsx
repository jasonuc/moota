import { useAuth } from "@/hooks/use-auth";
import { getUserUnactivatedPlantsCount } from "@/services/api/plants";
import { InfoIcon } from "lucide-react";
import { useEffect, useState } from "react";
import { Link } from "react-router";

export default function UnactivatedPlantsBanner() {
  const { user } = useAuth();
  const [unactivatedPlantsCount, setUnactivatedPlantsCount] = useState(0);

  useEffect(() => {
    getUserUnactivatedPlantsCount(user!.id).then(setUnactivatedPlantsCount);
  }, [user, setUnactivatedPlantsCount]);

  if (unactivatedPlantsCount <= 0) return null;
  return (
    <div className="w-fit flex items-center gap-x-1.5">
      <InfoIcon className="size-4" />
      <Link to="/plants/unactivated" className="underline text-blue-600">
        <small>
          You have {unactivatedPlantsCount} unactivated plant
          {(unactivatedPlantsCount ?? 0) > 1 ? "s" : ""}
        </small>
      </Link>
    </div>
  );
}
