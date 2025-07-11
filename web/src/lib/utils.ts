import { Coordinates } from "@/types/coordinates";
import { clsx, type ClassValue } from "clsx";
import { formatDate, formatDistanceToNow, isValid, parseJSON } from "date-fns";
import { twMerge } from "tailwind-merge";
import { EARTH_RADIUS_M } from "./constants";

export function cn(...inputs: ClassValue[]) {
  return twMerge(clsx(inputs));
}

export function formatDistance(meters: number) {
  if (meters >= 1000) {
    return `${(meters / 1000).toFixed(1)} km`;
  }
  return `${meters.toFixed(2)} m`;
}

export function formatHp(hp: number) {
  return `${hp.toFixed(0)}`;
}

export function startSentenceWithUppercase(sentence: string) {
  if (!sentence) return "";
  return sentence.charAt(0).toUpperCase() + sentence.slice(1);
}

export function formatPlantDate(dateString: string | undefined) {
  if (!dateString) return "Unknown";

  try {
    const date = parseJSON(dateString);
    return isValid(date) ? formatDate(date, "dd/MM/yy") : "Invalid date";
  } catch {
    return "Invalid date";
  }
}

export function formatRelativeTime(date: string | Date | number) {
  return formatDistanceToNow(date, { addSuffix: true });
}

export function formatCoordinates(lat: number, lon: number): string {
  const latDirection = lat >= 0 ? "N" : "S";
  const lonDirection = lon >= 0 ? "E" : "W";

  const latAbs = Math.abs(lat).toFixed(3);
  const lonAbs = Math.abs(lon).toFixed(3);

  return `${latAbs}°${latDirection}, ${lonAbs}°${lonDirection}`;
}

export function getDicebearThumbsUrl(seed?: string): string {
  if (!seed) {
    return `https://api.dicebear.com/9.x/thumbs/svg?seed=${"empty-seed"}&backgroundColor=transparent`;
  }
  return `https://api.dicebear.com/9.x/thumbs/svg?seed=${seed}&backgroundColor=transparent&shapeRotation=-20`;
}

export function getDicebearGlassUrl(seed?: string): string {
  if (!seed) {
    return `https://api.dicebear.com/9.x/glass/svg?seed=${"empty-seed"}`;
  }
  return `https://api.dicebear.com/9.x/glass/svg?seed=${seed}`;
}

export function haversineDistance(
  coord1: Coordinates,
  coord2: Coordinates
): number {
  const R = EARTH_RADIUS_M;
  const dLat = (coord2.Lat - coord1.Lat) * (Math.PI / 180);
  const dLon = (coord2.Lon - coord1.Lon) * (Math.PI / 180);

  const a =
    Math.sin(dLat / 2) * Math.sin(dLat / 2) +
    Math.cos(coord1.Lat * (Math.PI / 180)) *
      Math.cos(coord2.Lat * (Math.PI / 180)) *
      Math.sin(dLon / 2) *
      Math.sin(dLon / 2);
  const c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1 - a));

  return R * c;
}
