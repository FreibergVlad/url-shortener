import { BrowserRouter, Navigate, Route, Routes } from 'react-router-dom'
import Login from '@/pages/login'
import CreateAccount from '@/pages/create-account'
import CreateOrganization from '@/pages/create-organization'
import Links from '@/pages/links'
import Layout from '@/layouts/layout'

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
        <Route element={<Layout withSidebar={false} />}>
          <Route path="/organizations/create" element={<CreateOrganization withSidebar={false} />} />
          <Route path="*" element={<Navigate to="/organizations/create" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}

export function AuthenticatedWithOrganizationRouter() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<Layout withSidebar={true} />}>
          <Route index element={<Navigate to="/links" replace />} />
          <Route path="/links" element={<Links />} />
          <Route path="/organizations/create" element={<CreateOrganization withSidebar={true} />} />
          <Route path="*" element={<Navigate to="/links" replace />} />
        </Route>
      </Routes>
    </BrowserRouter>
  )
}
