import { useGeolocation } from "@/hooks/use-geolocation";
import { useEffect } from "react";
import { useNavigate } from "react-router";

type GeolocationAccuracyIndicatorProps = {
  children: React.ReactNode;
};

export function GeolocationAccuracyIndicator({
  children,
}: GeolocationAccuracyIndicatorProps) {
  const { withinAllowance } = useGeolocation();
  const navigation = useNavigate();

  useEffect(() => {
    if (!withinAllowance) {
      navigation("/low-geolocation-accuracy");
    }
  }, [withinAllowance, navigation]);

  if (withinAllowance === null) return null;
  if (!withinAllowance) {
    return null;
  } else {
    return children;
  }
}
