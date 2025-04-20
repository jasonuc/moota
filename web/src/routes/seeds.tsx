import Header from '@/components/header'
import { Button } from '@/components/ui/button'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/seeds')({
  component: RouteComponent,
})

function RouteComponent() {
  const seeds = [
    {
      id: 1,
      botanicalName: "Neoregalia Fosteriana",
      createdAt: "2023-10-01",
      count: 3,
    },
    {
      id: 2,
      botanicalName: "Spathiphyllum Wallisii",
      createdAt: "2023-10-02",
      count: 1,
    },
    {
      id: 3,
      botanicalName: "Ficus elastica",
      createdAt: "2023-10-03",
      count: 1,
    },
    {
      id: 4,
      botanicalName: "Monstera deliciosa",
      createdAt: "2023-10-04",
      count: 2,
    },
    {
      id: 5,
      botanicalName: "Aloe Vera",
      createdAt: "2023-10-05",
      count: 1,
    },
    {
      id: 6,
      botanicalName: "Cactus",
      createdAt: "2023-10-06",
      count: 2,
    },
  ]

  return (
    <div className="flex flex-col space-y-5 pb-10">
      <Header seedCount={10} />

      <h1 className="text-3xl font-heading mb-5">My Seeds</h1>

      <div className="grid grid-cols-3 md:grid-cols-4 gap-5">
        {seeds.map(({ id, botanicalName, count }) => (
          <Button asChild className="relative h-36" key={id}>
            <div className="size-full">
              {count > 1 && <small className="absolute right-2 -top-2 bg-background px-2 rounded-full">x{count}</small>}
              <p className="text-wrap text-center">{botanicalName}</p>
            </div>
          </Button>
        ))}
      </div>
    </div>
  )
}
