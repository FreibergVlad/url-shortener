import { useNavigate } from 'react-router-dom'
import { useState } from 'react'
import { cn } from '@/lib/shadcn-utils'
import { Button } from '@/components/shadcn/ui/button'
import { Input } from '@/components/shadcn/ui/input'
import { Label } from '@/components/shadcn/ui/label'
import Loader from '@/components/loader'
import { ErrorReason, ensureAppError } from '@/services/errors'
import { useAPIContext } from '@/contexts/api'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/shadcn/ui/card'

const ORGANIZATION_ALREADY_EXISTS_ERR_MSG = 'Organization with such slug already exists.'

const FORM_TITLE = 'Create Organization'
const FORM_DESCRIPTION = 'An organization groups your links and settings together, allowing you to easily manage and share them. Once created, you\'ll automatically become the owner and can invite other members later.'

const slugify = (s: string): string => {
  return s
    .replace(/^\s+|\s+$/g, '') // trim leading/trailing white space
    .toLowerCase()
    .replace(/[^a-z0-9 -]/g, '') // remove any non-alphanumeric characters
    .replace(/\s+/g, '-') // replace spaces with hyphens
    .replace(/-+/g, '-') // remove consecutive hyphens
}

function CreateOrganizationForm() {
  const navigate = useNavigate()
  const { useCreateOrganization } = useAPIContext()

  const [name, setName] = useState<string>('')
  const [slug, setSlug] = useState<string>('')
  const [slugManuallyChanged, setSlugManuallyChanged] = useState<boolean>(false)

  const [nameError, setNameError] = useState<string | null>(null)
  const [slugError, setSlugError] = useState<string | null>(null)
  const [generalError, setGeneralError] = useState<string | null>(null)

  const isLoading = useCreateOrganization.isPending

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault()
    setNameError(null)
    setSlugError(null)
    setGeneralError(null)

    if (!name.trim() || !slug.trim()) {
      return
    }

    try {
      await useCreateOrganization.mutateAsync({ name: name.trim(), slug: slug.trim() })
      await navigate('/')
    }
    catch (e: unknown) {
      const error = ensureAppError(e)
      if (error.reason === ErrorReason.BAD_REQUEST) {
        for (const detail of error.details) {
          if (detail.field === 'name') {
            setNameError(detail.description)
          }
          else if (detail.field === 'slug') {
            setSlugError(detail.description)
          }
        }
      }
      else if (error.reason === ErrorReason.ALREADY_EXISTS) {
        setSlugError(ORGANIZATION_ALREADY_EXISTS_ERR_MSG)
      }
      else {
        setGeneralError(error.friendlyMessage)
      }
    }
  }

  const onNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value)
    if (!slugManuallyChanged) {
      setSlug(slugify(e.target.value))
    }
  }

  const onSlugChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setSlug(e.target.value)
    setSlugManuallyChanged(true)
  }

  return (
    <>
      <form onSubmit={handleSubmit}>
        <div className="flex flex-col gap-6">
          <div className="grid gap-2">
            <Label htmlFor="name">Name</Label>
            <Input
              id="name"
              type="text"
              required
              onChange={onNameChange}
            />
            <p className="text-xs text-muted-foreground">
              Choose a descriptive and memorable name for your organization.
            </p>
            {nameError && <p className="text-red-500 text-sm text-left">{nameError}</p>}
          </div>
          <div className="grid gap-2">
            <Label htmlFor="slug">Slug</Label>
            <Input
              id="slug"
              type="text"
              required
              onChange={onSlugChange}
              value={slug}
            />
            <p className="text-xs text-muted-foreground">
              This will appear in your organization&apos;s URLs. Use only letters and numbers (no spaces or special characters).
            </p>
            {slugError && <p className="text-red-500 text-sm text-left">{slugError}</p>}
          </div>
          {generalError && <p className="text-red-500 text-sm text-center">{generalError}</p>}
          <Button type="submit" className="w-full">
            {isLoading ? (<><Loader className="mr-2 h-4 w-4 animate-spin" /></>) : ('Create')}
          </Button>
        </div>
      </form>
    </>
  )
}

export default function CreateOrganization({ className, withSidebar, ...props }: { withSidebar: boolean } & React.ComponentPropsWithoutRef<'div'>) {
  if (!withSidebar) {
    return (
      <div className="flex flex-row min-h-full justify-center items-center">
        <div className="w-full max-w-lg">
          <div className={cn('flex flex-col', className)} {...props}>
            <Card>
              <CardHeader>
                <CardTitle className="text-2xl">{FORM_TITLE}</CardTitle>
                <CardDescription>{FORM_DESCRIPTION}</CardDescription>
              </CardHeader>
              <CardContent>
                <CreateOrganizationForm />
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    )
  }

  return (
    <div className="flex flex-col p-8">
      <div className="w-full max-w-2xl space-y-6">
        <div className={cn('flex flex-col gap-6 justify-center', className)} {...props}>
          <div>
            <h2 className="text-2xl font-semibold">{FORM_TITLE}</h2>
            <p className="mt-2 text-sm text-muted-foreground">{FORM_DESCRIPTION}</p>
          </div>
          <CreateOrganizationForm />
        </div>
      </div>
    </div>
  )
};
