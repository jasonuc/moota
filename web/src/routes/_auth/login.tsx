import { createFileRoute, Link } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import LoginForm from '@/components/login-form'

export const Route = createFileRoute('/_auth/login')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <Card className='w-md h-fit'>
      <CardHeader className='text-center'>
        <CardTitle className='font-heading text-xl'>Login to your account</CardTitle>
        <CardDescription className='text-sm font-base'>
          Already have an account?{' '}
          <Link to='/register' className='text-blue-500 hover:underline'>
            Register
          </Link>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <LoginForm />
      </CardContent>
    </Card>
  )
}
