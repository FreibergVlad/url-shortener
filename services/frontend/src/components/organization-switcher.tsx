import { ChevronsUpDown, Plus } from 'lucide-react'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/shadcn/ui/dropdown-menu'

import {
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
} from '@/components/shadcn/ui/sidebar'
import { useSidebar } from '@/contexts/sidebar'
import { OrganizationMembership } from '@/models/organization'
import { Link } from 'react-router-dom'
import { AvatarIcon } from './avatar'

export const OrganizationSwitcher = ({
  organizationMemberships,
  currentOrganizationMembership,
  onOrganizationMembershipSelect,
}: {
  organizationMemberships: OrganizationMembership[]
  currentOrganizationMembership: OrganizationMembership
  onOrganizationMembershipSelect: (membership: OrganizationMembership) => void
}) => {
  const { isMobile } = useSidebar()

  return (
    <SidebarMenu>
      <SidebarMenuItem className="outline-none">
        <DropdownMenu>
          <DropdownMenuTrigger asChild className="outline-none ring-0 focus:ring-0 focus-visible:ring-0">
            <SidebarMenuButton
              size="lg"
              className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground cursor-pointer"
            >
              <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                <AvatarIcon name={currentOrganizationMembership.organization.slug} />
              </div>
              <div className="grid flex-1 text-left text-sm leading-tight">
                <span className="truncate font-semibold">
                  {currentOrganizationMembership.organization.slug}
                </span>
              </div>
              <ChevronsUpDown className="ml-auto" />
            </SidebarMenuButton>
          </DropdownMenuTrigger>
          <DropdownMenuContent
            className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
            align="start"
            side={isMobile ? 'bottom' : 'right'}
            sideOffset={4}
          >
            <DropdownMenuLabel className="text-xs text-muted-foreground">
              Organizations
            </DropdownMenuLabel>
            {organizationMemberships.map(membership => (
              <DropdownMenuItem
                key={membership.organization.slug}
                onClick={() => { onOrganizationMembershipSelect(membership) }}
                className="gap-2 p-2"
              >
                <div className="flex aspect-square size-8 items-center justify-center rounded-lg bg-sidebar-primary text-sidebar-primary-foreground">
                  <AvatarIcon name={membership.organization.slug} />
                </div>
                {membership.organization.slug}
              </DropdownMenuItem>
            ))}
            <DropdownMenuSeparator />
            <Link to="/organizations/create">
              <DropdownMenuItem className="gap-2 p-2">
                <div className="flex size-6 items-center justify-center rounded-md border bg-background">
                  <Plus className="size-4" />
                </div>
                <div className="font-medium text-muted-foreground">
                  Create
                </div>
              </DropdownMenuItem>
            </Link>
          </DropdownMenuContent>
        </DropdownMenu>
      </SidebarMenuItem>
    </SidebarMenu>
  )
}
