import { Stats } from "@/types/user";
import { useEffect, useState } from "react";
import useWebSocket from "react-use-websocket";

export function useStats() {
  const WS_URL = "/api/stats";
  const [stats, setStats] = useState<Stats>({
    plant: { alive: 0, deceased: 0 },
    seed: { Planted: 0, unused: 0 },
  });
  const { lastJsonMessage } = useWebSocket<Stats | null>(WS_URL);

  useEffect(() => {
    if (!lastJsonMessage) return;
    setStats(lastJsonMessage);
    console.log(lastJsonMessage);
  }, [lastJsonMessage]);

  return { stats };
}
