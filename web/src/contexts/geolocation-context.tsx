import { GeolocationState } from "@uidotdev/usehooks";
import { createContext } from "react";

interface GeolocationContextState extends GeolocationState {
  withinAllowance: boolean;
}

export const GeolocationContext = createContext<GeolocationContextState>({
  loading: true,
  accuracy: null,
  altitude: null,
  altitudeAccuracy: null,
  error: null,
  latitude: null,
  longitude: null,
  heading: null,
  speed: null,
  timestamp: null,
  withinAllowance: false,
});
