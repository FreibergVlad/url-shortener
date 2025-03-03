import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import Login from "./pages/login";
import CreateAccount from "./pages/create-account";
import Loader from "./components/loader";
import { useUserContext } from "./contexts/user";
import CreateOrganization from "./pages/create-organization";
import Dashboard from "./pages/dashboard";

export default function App() {
  const {user, organizationMemberships, isLoading} = useUserContext();

  if (isLoading) {
    return <Loader />;
  }

  if (!user) {
    return (
      <BrowserRouter>
        <Routes>
          <Route path="/login" element={<Login />} />
          <Route path="/create-account" element={<CreateAccount />} />
          <Route path="*" element={<Navigate to="/login" replace />} />
        </Routes>
      </BrowserRouter>
    );
  }

  if (!organizationMemberships || !organizationMemberships.length) {
    return (
      <BrowserRouter>
        <Routes>
          <Route path="/organizations/create" element={<CreateOrganization />} />
          <Route path="*" element={<Navigate to="/organizations/create" replace />} />
        </Routes>
      </BrowserRouter>
    );
  }

  return (
    <BrowserRouter>
      <Routes>
        <Route path="/dashboard" element={<Dashboard />} />
        <Route path="/organizations/create" element={<CreateOrganization />} />
        <Route path="*" element={<Navigate to="/dashboard" replace />} />
      </Routes>
    </BrowserRouter>
  );
}