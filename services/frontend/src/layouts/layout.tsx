import { NavUser } from '@/components/nav-user'
import { Separator } from '@/components/shadcn/ui/separator'
import { SidebarInset } from '@/components/shadcn/ui/sidebar'
import { Sidebar } from '@/components/sidebar'
import { SidebarProvider } from '@/providers/sidebar'
import { Outlet } from 'react-router-dom'

export default function Layout({ withSidebar }: { withSidebar: boolean }) {
  return (
    <SidebarProvider>
      {withSidebar && <Sidebar />}
      <SidebarInset>
        <div className="flex flex-col min-h-screen">
          <header className="flex h-16 items-center justify-between px-4 gap-2 bg-sidebar">
            <div className="flex items-center gap-2 px-4">
              <Separator orientation="vertical" className="mr-2 h-4" />
            </div>
            <div className="flex items-center gap-2">
              <NavUser />
            </div>
          </header>
          <Separator orientation="horizontal" className="w-full bg-border" />
          <main className="flex-1 overflow-y-auto">
            <Outlet />
          </main>
        </div>
      </SidebarInset>
    </SidebarProvider>
  )
}
