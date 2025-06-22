import { GeolocationContext } from "@/contexts/geolocation-context";
import { useGeolocation } from "@uidotdev/usehooks";
import { useEffect, useState } from "react";

type GeolocationProviderProps = {
  children: React.ReactNode;
};

const DISTANCE_ACCURACY_ALLOWANCE = 10;

export default function GeolocationProvider({
  children,
}: GeolocationProviderProps) {
  const geolocation = useGeolocation({
    enableHighAccuracy: true,
    maximumAge: 0,
    timeout: 1000,
  });
  const [withinAllowance, setWithinAllowance] = useState<boolean>(false);

  useEffect(() => {
    if (geolocation.accuracy !== null) {
      setWithinAllowance(geolocation.accuracy <= DISTANCE_ACCURACY_ALLOWANCE);
    }
  }, [geolocation.accuracy]);

  return (
    <GeolocationContext.Provider value={{ ...geolocation, withinAllowance }}>
      {children}
    </GeolocationContext.Provider>
  );
}
