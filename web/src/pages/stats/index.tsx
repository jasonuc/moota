import { Stats } from "@/types/user";
import { useEffect, useState } from "react";

export default function StatsPage() {
  const [stockData, setStockData] = useState<Stats | null>(null);

  useEffect(() => {
    // opening a connection to the server to begin receiving events from it
    const eventSource = new EventSource("http://localhost:5173/api/stats");

    // attaching a handler to receive message events
    eventSource.onmessage = (event) => {
      const data = JSON.parse(event.data);
      console.log("stockData", data);
      setStockData(data);
    };

    // terminating the connection on component unmount
    return () => {
        eventSource.close(); // Cleanup on unmount
      };
  }, []);
  if (!stockData) {
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
