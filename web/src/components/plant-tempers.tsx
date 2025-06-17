import { Progress } from "./ui/progress";

type PlantTempersProps = {
  woe?: number;
  frolic?: number;
  malice?: number;
  dread?: number;
  loading?: boolean;
};

export default function PlantTempers({
  woe,
  dread,
  frolic,
  malice,
}: PlantTempersProps) {
  const items = [
    { label: "woe", value: woe },
    { label: "dread", value: dread },
    { label: "frolic", value: frolic },
    { label: "malice", value: malice },
  ];

  return (
    <div>
      <h3 className="font-heading mb-3">Tempers</h3>
      <div className="grid grid-cols-2 grid-rows-2 gap-y-5 gap-x-10 md:gap-10">
        {items.map(({ label, value }) => (
          <div key={label}>
            <p className="mb-1 ml-1 md:ml-0 md:mb-0 grow-0 italic pb-3 text-left">
              {label}
            </p>
            <Progress
              className="h-5 rounded-md md:rounded-base md:min-h-8"
              value={typeof value === "number" ? (value / 5) * 100 : 0}
            />
          </div>
        ))}
      </div>
    </div>
  );
}
