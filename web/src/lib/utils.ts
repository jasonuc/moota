import { clsx, type ClassValue } from "clsx";
import { twMerge } from "tailwind-merge";

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
