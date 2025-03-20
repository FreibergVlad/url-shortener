import { Link, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { cn } from '@/lib/shadcn-utils'
import { Button } from '@/components/shadcn/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/ui/card'
import { Input } from '@/components/shadcn/ui/input'
import { Label } from '@/components/shadcn/ui/label'
import Loader from '@/components/loader'
import { ensureAppError } from '@/services/errors'
import { useAuthContext } from '@/contexts/auth'

export default function Login({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
  const navigate = useNavigate()
  const { authenticateUser } = useAuthContext()

  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')

  const [generalError, setGeneralError] = useState<string | null>(null)

  const isLoading = authenticateUser.isPending

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()

    if (!email.trim() || !password.trim()) {
      return
    }

    try {
      await authenticateUser.mutateAsync({ email: email.trim(), password: password.trim() })
      await navigate('/')
    }
    catch (e: unknown) {
      const error = ensureAppError(e)
      setGeneralError(error.friendlyMessage)
    }
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <div className={cn('flex flex-col gap-6', className)} {...props}>
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl text-center">Login</CardTitle>
              <CardDescription className="text-center">
                Enter your email below to login to your account
              </CardDescription>
            </CardHeader>
            <CardContent>
              <form onSubmit={handleSubmit}>
                <div className="flex flex-col gap-6">
                  <div className="grid gap-2">
                    <Label htmlFor="email">Email</Label>
                    <Input
                      id="email"
                      type="email"
                      placeholder="m@example.com"
                      required
                      onChange={(e) => { setEmail(e.target.value) }}
                    />
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="password">Password</Label>
                    <Input
                      id="password"
                      type="password"
                      required
                      onChange={(e) => { setPassword(e.target.value) }}
                    />
                  </div>
                  {generalError && <p className="text-red-500 text-sm text-center">{generalError}</p>}
                  <Button type="submit" className="w-full">
                    {isLoading ? (<><Loader className="mr-2 h-4 w-4 animate-spin" /></>) : ('Login')}
                  </Button>
                </div>
                <div className="mt-4 text-center text-sm">
                  Don&apos;t have an account?
                  {' '}
                  <Link to="/create-account" className="underline underline-offset-4">
                    Sign up
                  </Link>
                </div>
              </form>
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  )
};
