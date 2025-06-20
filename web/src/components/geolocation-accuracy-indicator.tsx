import { useGeolocation } from "@uidotdev/usehooks";
import { useNavigate } from "react-router";

type GeolocationAccuracyIndicatorProps = {
  children: React.ReactNode;
};

export default function GeolocationAccuracyIndicator({
  children,
}: GeolocationAccuracyIndicatorProps) {
  const { accuracy } = useGeolocation();
  const navigation = useNavigate();

  if (!accuracy) return null;

  if (accuracy > 5) {
    navigation("/low-geolocation-accuracy");
    return null;
  } else {
    return children;
  }
}
