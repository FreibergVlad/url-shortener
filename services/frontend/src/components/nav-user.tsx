import {
  ChevronsUpDown,
  LogOut,
  UserPen,
} from 'lucide-react'

import { Avatar } from '@/components/shadcn/ui/avatar'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuGroup,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/shadcn/ui/dropdown-menu'
import { SidebarMenuButton } from '@/components/shadcn/ui/sidebar'
import { useUserContext } from '@/contexts/user'
import { useAuthContext } from '@/contexts/auth'
import { AvatarIcon } from './avatar'

export function NavUser() {
  const { user } = useUserContext()
  const { logout } = useAuthContext()

  if (!user) {
    return null
  }

  return (
    <DropdownMenu>
      <DropdownMenuTrigger asChild>
        <SidebarMenuButton
          size="lg"
          className="data-[state=open]:bg-sidebar-accent data-[state=open]:text-sidebar-accent-foreground outline-none ring-0 focus:ring-0 focus-visible:ring-0 cursor-pointer"
        >
          <Avatar className="h-8 w-8 rounded-lg">
            <AvatarIcon name={user.fullName} />
          </Avatar>
          <div className="grid flex-1 text-left text-sm leading-tight">
            <span className="truncate font-semibold">{user.fullName}</span>
            <span className="truncate text-xs">{user.email}</span>
          </div>
          <ChevronsUpDown className="ml-auto size-4" />
        </SidebarMenuButton>
      </DropdownMenuTrigger>
      <DropdownMenuContent
        className="w-[--radix-dropdown-menu-trigger-width] min-w-56 rounded-lg"
        align="end"
        sideOffset={4}
      >
        <DropdownMenuLabel className="p-0 font-normal outline-none">
          <div className="flex items-center gap-2 px-1 py-1.5 text-left text-sm">
            <Avatar className="h-8 w-8 rounded-lg">
              <AvatarIcon name={user.fullName} />
            </Avatar>
            <div className="grid flex-1 text-left text-sm leading-tight">
              <span className="truncate font-semibold">{user.fullName}</span>
              <span className="truncate text-xs">{user.email}</span>
            </div>
          </div>
        </DropdownMenuLabel>
        <DropdownMenuSeparator />
        <DropdownMenuGroup>
          <DropdownMenuItem>
            <UserPen />
            Account
          </DropdownMenuItem>
        </DropdownMenuGroup>
        <DropdownMenuSeparator />
        <DropdownMenuItem onClick={logout}>
          <LogOut />
          Log out
        </DropdownMenuItem>
      </DropdownMenuContent>
    </DropdownMenu>
  )
}
