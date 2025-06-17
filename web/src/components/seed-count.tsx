import { HexagonIcon } from "lucide-react";

type SeedCountProps = {
  number?: number;
  size: number;
};

export default function SeedCount({ number, size }: SeedCountProps) {
  return (
    <div className="relative inline-flex">
      <HexagonIcon size={size} strokeWidth={1} />
      <div className="absolute inset-0 flex items-center justify-center text-xs">
        {number}
      </div>
    </div>
  );
}
