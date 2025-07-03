import { useEffect, useState } from 'react';

function useSSE<T>(url: string) {
  const [data, setData] = useState<T | null>(null);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let eventSource = new EventSource(url);

    // Handle incoming data
    eventSource.onmessage = (event) => {
      const newData = JSON.parse(event.data);
      console.log("newData", newData);
      setData(newData);
    };

    // Handle errors
    eventSource.onerror = (error) => {
      console.error('SSE error:', error);
      setError("Error occurred while trying to fetch data");
      // Attempt to reconnect after a delay
      setTimeout(() => {
        eventSource = new EventSource(url);
      }, 2000); // Reconnect after 2 seconds
    };

    // Cleanup when component unmounts
    return () => eventSource.close();
  }, [url]);

  return { data, error };
}

export default useSSE;