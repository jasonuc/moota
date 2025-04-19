import DashboardHeader from '@/components/dashboard-header'
import { createFileRoute } from '@tanstack/react-router'

export const Route = createFileRoute('/seeds')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <div className='flex flex-col space-y-5 pb-10'>
        <DashboardHeader seedCount={10} />

        <h1 className='text-3xl font-heading mb-5'>My Seeds</h1>
    </div>
  )
}
