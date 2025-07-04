import { GeolocationContext } from "@/contexts/geolocation-context";
import { GEOLOCATION_INACURACCY_TOLERANCE } from "@/lib/constants";
import { useGeolocation } from "@uidotdev/usehooks";
import { LucideAlertTriangle } from "lucide-react";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";
import { toast } from "sonner";

type GeolocationProviderProps = {
  children: React.ReactNode;
};

export default function GeolocationProvider({
  children,
}: GeolocationProviderProps) {
  const geolocation = useGeolocation({
    enableHighAccuracy: true,
    maximumAge: 0,
    timeout: 1000,
  });
  const [withinAllowance, setWithinAllowance] = useState(false);

  const navigate = useNavigate();
  const { pathname } = useLocation();

  useEffect(() => {
    if (geolocation.error) {
      if (geolocation.error.code === geolocation.error.PERMISSION_DENIED) {
        const nonGeolocationDepentendPaths = [
          "/",
          "/settings",
          "/plants/graveyard",
        ];
        if (!nonGeolocationDepentendPaths.includes(pathname)) {
          navigate("/geolocation-dissallowed");
        }
      }
    }
  }, [geolocation.error, pathname, navigate]);

  useEffect(() => {
    if (geolocation.accuracy !== null) {
      setWithinAllowance(
        geolocation.accuracy <= GEOLOCATION_INACURACCY_TOLERANCE
      );
    }
  }, [geolocation.accuracy]);

  useEffect(() => {
    if (geolocation.accuracy === null && geolocation.error !== null) {
      if (geolocation.error.code === geolocation.error.POSITION_UNAVAILABLE) {
        toast.dismiss();
        toast("Location unavailable", {
          description:
            "Application will remain available but various features will be disabled until location is available again.",
          icon: <LucideAlertTriangle />,
          duration: 2500,
          dismissible: true,
        });
      }
    }
  }, [geolocation.accuracy, geolocation.error]);

  return (
    <GeolocationContext.Provider value={{ ...geolocation, withinAllowance }}>
      {children}
    </GeolocationContext.Provider>
  );
}
