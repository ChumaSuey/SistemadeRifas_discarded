import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { CssBaseline } from '@mui/material';
import { AuthProvider } from './context/AuthContext';
import { SnackbarProvider } from './context/SnackbarContext';
import { AppThemeProvider } from './context/ThemeContext';
import Layout from './components/Layout';
import ProtectedRoute from './components/ProtectedRoute';
import HomePage from './pages/HomePage';
import LoginPage from './pages/LoginPage';
import RaffleDetailPage from './pages/RaffleDetailPage';
import DashboardPage from './pages/DashboardPage';
import ManageRafflePage from './pages/ManageRafflePage';
import LiveDrawPage from './pages/LiveDrawPage';
import VerifierPage from './pages/VerifierPage';
import MyTicketsPage from './pages/MyTicketsPage';
import AdminPage from './pages/AdminPage';

export default function App() {
  return (
    <AppThemeProvider>
      <CssBaseline />
      <BrowserRouter>
        <AuthProvider>
          <SnackbarProvider>
            <Routes>
              <Route element={<Layout />}>
                <Route path="/" element={<HomePage />} />
                <Route path="/login" element={<LoginPage />} />
                <Route path="/raffles/:id" element={<RaffleDetailPage />} />

                <Route path="/dashboard" element={
                  <ProtectedRoute roles={['raffler', 'admin']}><DashboardPage /></ProtectedRoute>
                } />
                <Route path="/raffles/:id/manage" element={
                  <ProtectedRoute roles={['raffler', 'admin']}><ManageRafflePage /></ProtectedRoute>
                } />
                <Route path="/raffles/:id/live" element={
                  <ProtectedRoute roles={['raffler', 'admin']}><LiveDrawPage /></ProtectedRoute>
                } />
                <Route path="/verifier" element={
                  <ProtectedRoute roles={['verifier', 'admin']}><VerifierPage /></ProtectedRoute>
                } />
                <Route path="/tickets/my" element={
                  <ProtectedRoute><MyTicketsPage /></ProtectedRoute>
                } />
                <Route path="/admin" element={
                  <ProtectedRoute roles={['admin']}><AdminPage /></ProtectedRoute>
                } />
              </Route>
            </Routes>
          </SnackbarProvider>
        </AuthProvider>
      </BrowserRouter>
    </AppThemeProvider>
  );
}
