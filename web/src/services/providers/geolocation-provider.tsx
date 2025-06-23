import { GeolocationContext } from "@/contexts/geolocation-context";
import { GEOLOCATION_DISTANCE_ACCURACY_ALLOWANCE } from "@/lib/constants";
import { useGeolocation } from "@uidotdev/usehooks";
import { useEffect, useState } from "react";
import { useLocation, useNavigate } from "react-router";

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
  const [withinAllowance, setWithinAllowance] = useState<boolean>(false);

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
        geolocation.accuracy <= GEOLOCATION_DISTANCE_ACCURACY_ALLOWANCE
      );
    }
  }, [geolocation.accuracy]);

  return (
    <GeolocationContext.Provider value={{ ...geolocation, withinAllowance }}>
      {children}
    </GeolocationContext.Provider>
  );
}
