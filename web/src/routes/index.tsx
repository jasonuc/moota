import { Button } from '@/components/ui/button';
import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute("/")({
  component: Index,
})

function Index() {
  return (
    <div>
      <div className='flex w-full justify-between items-center'>
        <div className='flex items-center justify-center space-x-2'>
          <img src='./moota.png' width={45} height={45} />
          <h1 className='font-bold text-3xl'>Moota</h1>
        </div>

        <div>
          <Link to='/dashboard'>
            <Button>
              Home
            </Button>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default Index;