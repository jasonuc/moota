import { LogoWithText } from '@/components/logo';
import { Button } from '@/components/ui/button';
import { createFileRoute, Link } from '@tanstack/react-router'

export const Route = createFileRoute("/")({
  component: Index,
})

function Index() {
  return (
    <div className='flex flex-col space-y-5 grow'>
      <div className='flex w-full justify-between items-center'>
        <Link to='/'>
          <LogoWithText />
        </Link>

        <div className='flex space-x-2'>
          <Link to='/dashboard'>
            <Button>
              Home
            </Button>
          </Link>
          
          <Link to='/register'>
            <Button>
              Register
            </Button>
          </Link>
          
          <Link to='/login'>
            <Button>
              Login
            </Button>
          </Link>
        </div>
      </div>
    </div>
  )
}

export default Index;