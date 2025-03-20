import { Link, useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/ui/card'
import { cn } from '@/lib/shadcn-utils'
import { Button } from '@/components/shadcn/ui/button'
import { Input } from '@/components/shadcn/ui/input'
import { Label } from '@/components/shadcn/ui/label'
import Loader from '@/components/loader'
import { ErrorReason, ensureAppError } from '@/services/errors'
import { useAuthContext } from '@/contexts/auth'

const EMAIL_ALREADY_EXISTS = 'User with such email already exists.'

export default function CreateAccount({ className, ...props }: React.ComponentPropsWithoutRef<'div'>) {
  const navigate = useNavigate()
  const { createUser, authenticateUser } = useAuthContext()

  const [email, setEmail] = useState<string>('')
  const [password, setPassword] = useState<string>('')

  const [emailError, setEmailError] = useState<string | null>(null)
  const [passwordError, setPasswordError] = useState<string | null>(null)
  const [generalError, setGeneralError] = useState<string | null>(null)

  const isLoading = createUser.isPending || authenticateUser.isPending

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setEmailError(null)
    setPasswordError(null)
    setGeneralError(null)

    if (!email.trim() || !password.trim()) {
      return
    }
    const credentials = { email: email.trim(), password: password.trim() }

    try {
      await createUser.mutateAsync(credentials)
      await authenticateUser.mutateAsync(credentials)

      await navigate('/')
    }
    catch (e: unknown) {
      const error = ensureAppError(e)
      if (error.reason === ErrorReason.BAD_REQUEST) {
        for (const detail of error.details) {
          if (detail.field === 'email') {
            setEmailError(detail.description)
          }
          else if (detail.field === 'password') {
            setPasswordError(detail.description)
          }
        }
      }
      else if (error.reason === ErrorReason.ALREADY_EXISTS) {
        setEmailError(EMAIL_ALREADY_EXISTS)
      }
      else {
        setGeneralError(error.friendlyMessage)
      }
    }
  }

  return (
    <div className="flex min-h-svh w-full items-center justify-center p-6 md:p-10">
      <div className="w-full max-w-sm">
        <div className={cn('flex flex-col gap-6', className)} {...props}>
          <Card>
            <CardHeader>
              <CardTitle className="text-2xl text-center">Sign Up</CardTitle>
              <CardDescription className="text-center">
                Enter your details below to create an account
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
                    {emailError && <p className="text-red-500 text-sm text-left">{emailError}</p>}
                  </div>
                  <div className="grid gap-2">
                    <Label htmlFor="password">
                      Password
                    </Label>
                    <Input
                      id="password"
                      type="password"
                      required
                      onChange={(e) => { setPassword(e.target.value) }}
                    />
                    {passwordError && <p className="text-red-500 text-sm text-left">{passwordError}</p>}
                  </div>
                  {generalError && <p className="text-red-500 text-sm text-center">{generalError}</p>}
                  <Button type="submit" className="w-full">
                    {isLoading ? (<><Loader className="mr-2 h-4 w-4 animate-spin" /></>) : ('Sign Up')}
                  </Button>
                </div>
                <div className="mt-4 text-center text-sm">
                  Already have an account?
                  {' '}
                  <Link to="/login" className="underline underline-offset-4">
                    Login
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
