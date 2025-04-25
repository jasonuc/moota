import { createFileRoute, Link } from '@tanstack/react-router'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import RegisterForm from '@/components/register-form'

export const Route = createFileRoute('/_auth/register')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <Card className='w-md h-fit'>
      <CardHeader className='text-center'>
        <CardTitle className='font-heading text-xl'>Create an account</CardTitle>
        <CardDescription className='text-sm font-base'>
          Already have an account?{' '}
          <Link to='/login' className='text-blue-500 underline-offset-2 underline'>
            Login
          </Link>
        </CardDescription>
      </CardHeader>
      <CardContent>
        <RegisterForm />
      </CardContent>
    </Card>
  )
}
