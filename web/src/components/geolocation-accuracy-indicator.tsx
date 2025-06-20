import { useGeolocation } from "@uidotdev/usehooks";
import { useNavigate } from "react-router";
import ProtectedRoute from "./protected";

type GeolocationAccuracyIndicatorProps = {
  children: React.ReactNode;
};

type ProtectedGeolocationAccuracyIndicatorRouteProps = {
  children: React.ReactNode;
};

export function GeolocationAccuracyIndicator({
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

export const ProtectedGeolocationAccuracyIndicatorRoute = ({
  children,
}: ProtectedGeolocationAccuracyIndicatorRouteProps) => (
  <ProtectedRoute>
    <GeolocationAccuracyIndicator>{children}</GeolocationAccuracyIndicator>
  </ProtectedRoute>
);
