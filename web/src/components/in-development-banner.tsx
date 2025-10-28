import { useLocalStorage } from "@uidotdev/usehooks";
import { AlertCircleIcon, XIcon } from "lucide-react";

export default function InDevelopmentBanner() {
  const [showInDevelopmentBanner, setShowInDevelopmentBanner] = useLocalStorage(
    "show-in-development-banner",
    true,
  );

  if (!showInDevelopmentBanner) return null;

  return (
    <div className="p-2 flex items-center justify-between gap-x-1.5 bg-amber-600/70">
      <div className="flex items-center gap-x-1.5">
        <AlertCircleIcon className="size-3" />
        <small className="font-semibold">
          Currently in development. Breaking changes could occur.
        </small>
      </div>

      <button
        onClick={() => {
          setShowInDevelopmentBanner(false);
        }}
      >
        <div>
          <XIcon className="size-4" />
        </div>
      </button>
    </div>
  );
}
