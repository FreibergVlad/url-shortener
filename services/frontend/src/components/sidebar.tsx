import { Link, Settings, Earth } from 'lucide-react'

import {
  Sidebar as BaseSidebar,
  SidebarContent,
  SidebarGroup,
  SidebarGroupContent,
  SidebarHeader,
  SidebarMenu,
  SidebarMenuButton,
  SidebarMenuItem,
  SidebarRail,
} from '@/components/shadcn/ui/sidebar'
import { OrganizationSwitcher } from './organization-switcher'
import { useUserContext } from '@/contexts/user'

const items = [
  {
    title: 'Links',
    url: '/links',
    icon: Link,
  },
  {
    title: 'Domains',
    url: '#',
    icon: Earth,
  },
  {
    title: 'Settings',
    url: '#',
    icon: Settings,
  },
]

export function Sidebar() {
  const { currentOrganizationMembership, organizationMemberships, setCurrentOrganizationMembership } = useUserContext()

  return (
    <BaseSidebar collapsible="icon">
      <SidebarHeader>
        {currentOrganizationMembership && organizationMemberships && (
          <OrganizationSwitcher
            organizationMemberships={organizationMemberships}
            currentOrganizationMembership={currentOrganizationMembership}
            onOrganizationMembershipSelect={setCurrentOrganizationMembership}
          />
        )}
      </SidebarHeader>
      <SidebarContent>
        <SidebarGroup>
          <SidebarGroupContent>
            <SidebarMenu>
              {items.map(item => (
                <SidebarMenuItem key={item.title}>
                  <SidebarMenuButton asChild>
                    <a href={item.url}>
                      <item.icon />
                      <span>{item.title}</span>
                    </a>
                  </SidebarMenuButton>
                </SidebarMenuItem>
              ))}
            </SidebarMenu>
          </SidebarGroupContent>
        </SidebarGroup>
      </SidebarContent>
      <SidebarRail />
    </BaseSidebar>
  )
}
