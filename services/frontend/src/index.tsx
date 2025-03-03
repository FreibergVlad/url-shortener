import { StrictMode } from 'react'
import { createRoot } from 'react-dom/client'
import './index.css'
import App from './App.tsx'
import { AuthProvider } from './contexts/auth.tsx'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { APIProvider } from './contexts/api.tsx'
import { UserProvider } from './contexts/user.tsx'

const queryClient = new QueryClient();

createRoot(document.getElementById('root')!).render(
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
