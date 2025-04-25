import { createFileRoute } from '@tanstack/react-router'
import RegisterForm from '@/components/register-form'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'

export const Route = createFileRoute('/_auth/register')({
  component: RouteComponent,
})

function RouteComponent() {
  return (
    <Card className='w-md h-fit'>
      <CardHeader className='text-center'>
        <CardTitle className='font-heading text-xl'>Create an account</CardTitle>
      </CardHeader>
      <CardContent>
        <RegisterForm />
      </CardContent>
    </Card>
  )
}
