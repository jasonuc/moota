import { GeolocationContext } from "@/contexts/geolocation-context";
import { useContext } from "react";

export function useGeolocation() {
  return useContext(GeolocationContext);
}
