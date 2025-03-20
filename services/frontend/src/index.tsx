import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import { AuthProvider } from './providers/auth.tsx'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { APIProvider } from './providers/api.tsx'
import { UserProvider } from './providers/user.tsx'
import App from './app/index.tsx'

const queryClient = new QueryClient()

const rootElement = document.getElementById('root')
if (!rootElement) {
  throw new Error('Root element not found')
}

createRoot(rootElement).render(
  <StrictMode>
    <QueryClientProvider client={queryClient}>
      <AuthProvider>
        <APIProvider>
          <UserProvider>
            <App />
          </UserProvider>
        </APIProvider>
      </AuthProvider>
    </QueryClientProvider>
  </StrictMode>,
)
