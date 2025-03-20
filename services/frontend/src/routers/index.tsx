import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from '@/pages/login'
import CreateAccount from '@/pages/create-account'
import CreateOrganization from '@/pages/create-organization'
import Links from '@/pages/links'

export function UnauthenticatedRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route path="/create-account" element={<CreateAccount />} />
        <Route path="*" element={<Navigate to="/login" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export function AuthenticatedWithoutOrganizationRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/organizations/create" element={<CreateOrganization />} />
        <Route path="*" element={<Navigate to="/organizations/create" replace />} />
      </Routes>
    </BrowserRouter>
  )
}

export function AuthenticatedWithOrganizationRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/links" element={<Links />} />
        <Route path="/organizations/create" element={<CreateOrganization />} />
        <Route path="*" element={<Navigate to="/links" replace />} />
      </Routes>
    </BrowserRouter>
  )
}
