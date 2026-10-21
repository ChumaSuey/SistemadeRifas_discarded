import { useState } from 'react';
import { Container, TextField, Button, Typography, Box, Alert, Tabs, Tab, CircularProgress } from '@mui/material';
import { useNavigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

export default function LoginPage() {
  const [tab, setTab] = useState(0);
  const [name, setName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');
  const { login, register, loading } = useAuth();
  const navigate = useNavigate();
  const location = useLocation();

  const from = (location.state as { from?: string })?.from || '/dashboard';

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');
    try {
      if (tab === 0) {
        await login(email, password);
      } else {
        await register(name, email, password);
      }
      navigate(from, { replace: true });
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { detail?: string; message?: string } } })?.response?.data?.detail
        || (err as { response?: { data?: { message?: string } } })?.response?.data?.message
        || 'Error al iniciar sesión';
      setError(msg);
    }
  };

  return (
    <Container maxWidth="xs" sx={{ mt: { xs: 4, md: 8 } }}>
      <Typography variant="h4" align="center" gutterBottom>Sistema de Rifas</Typography>

      <Tabs value={tab} onChange={(_, v) => setTab(v)} centered sx={{ mb: 3 }}>
        <Tab label="Entrar" />
        <Tab label="Registrarse" />
      </Tabs>

      {error && <Alert severity="error" sx={{ mb: 2 }}>{error}</Alert>}

      <Box component="form" onSubmit={handleSubmit}>
        {tab === 1 && (
          <TextField fullWidth label="Nombre" value={name} onChange={(e) => setName(e.target.value)} margin="normal" required />
        )}
        <TextField fullWidth label="Email" type="email" value={email} onChange={(e) => setEmail(e.target.value)} margin="normal" required />
        <TextField fullWidth label="Contraseña" type="password" value={password} onChange={(e) => setPassword(e.target.value)} margin="normal" required />
        <Button type="submit" variant="contained" fullWidth size="large" disabled={loading} sx={{ mt: 2 }}>
          {loading ? <CircularProgress size={24} color="inherit" /> : tab === 0 ? 'Entrar' : 'Registrarse'}
        </Button>
      </Box>
    </Container>
  );
}
