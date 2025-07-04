import { Stats } from "@/types/user";
import {
  useEventSource,
  useEventSourceListener,
} from "@react-nano/use-event-source";
import { useState } from "react";

export default function StatsPage() {
  const [stockData, setMesages] = useState<Stats | null>(null);

  const [eventSource, eventSourceStatus] = useEventSource("api/stats", false);
  useEventSourceListener(eventSource, ["data"], (evt) => {
    setMesages(JSON.parse(evt.data));
  });
  if (eventSourceStatus === "error") {
    return <div>EventSource error</div>;
  }
  if (eventSourceStatus === "init") {
    return <div>Loading...</div>;
  }
  return (
    <div className="flex flex-col space-y-5 pb-10 w-full">
      <h1 className="text-3xl font-heading mb-5">Stats</h1>
      <h1 className="text-6xl md:text-8xl font-heading text-foreground mb-6 leading-tight">
        Plant Your
        <br />
        <span className="text-main relative">
          Adventure
          <div className="absolute -bottom-2 left-1/2 transform -translate-x-1/2 w-32 h-1 bg-main rounded-full" />
        </span>
      </h1>

      <p className="text-xl md:text-2xl text-foreground mb-12 max-w-4xl mx-auto leading-relaxed font-medium opacity-80">
        Turn the real world into your playground. Plant virtual seeds anywhere,
        return to care for them, and transform daily routines into epic gaming
        adventures.
      </p>
      <p>Plant Count: {stockData?.plant?.alive}</p>
    </div>
  );
}
