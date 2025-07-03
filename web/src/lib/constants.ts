export const EARTH_RADIUS_M = 6.378e6;
export const GEOLOCATION_INACURACCY_TOLERANCE = import.meta.env.PROD
  ? 20
  : import.meta.env.GEOLOCATION_INACURACCY_TOLERANCE ?? 50; // TODO: This value is still experimental

export const SERVER_404_MESSAGE = "Requested resource could not be found";
