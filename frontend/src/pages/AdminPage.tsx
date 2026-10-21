import { useState, useEffect, useCallback } from 'react';
import { Typography, Box, Button, Table, TableBody, TableCell, TableHead, TableRow, Chip, Paper, Grid, CircularProgress, Dialog, DialogTitle, DialogContent, DialogActions, Pagination } from '@mui/material';
import PeopleIcon from '@mui/icons-material/People';
import ConfirmationNumberIcon from '@mui/icons-material/ConfirmationNumber';
import CasinoIcon from '@mui/icons-material/Casino';
import GavelIcon from '@mui/icons-material/Gavel';
import api from '../api/client';
import type { User, PaginatedResponse } from '../types';
import { roleLabel } from '../utils/formatters';
import { useSnackbar } from '../context/SnackbarContext';

interface AdminStats {
  total_users: number;
  total_raffles: number;
  total_tickets: number;
  pending_tickets: number;
}

function StatCard({ icon, label, value }: { icon: React.ReactNode; label: string; value: number }) {
  return (
    <Paper sx={{ p: 2, textAlign: 'center', height: '100%' }} variant="outlined">
      <Box sx={{ color: 'primary.main', mb: 1 }}>{icon}</Box>
      <Typography variant="h4" sx={{ fontWeight: 700 }}>{value}</Typography>
      <Typography variant="body2" color="text.secondary">{label}</Typography>
    </Paper>
  );
}

const USERS_PER_PAGE = 10;

export default function AdminPage() {
  const [users, setUsers] = useState<User[]>([]);
  const [stats, setStats] = useState<AdminStats>({ total_users: 0, total_raffles: 0, total_tickets: 0, pending_tickets: 0 });
  const [loading, setLoading] = useState(true);
  const [suspendTarget, setSuspendTarget] = useState<User | null>(null);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const { showSnackbar } = useSnackbar();

  const fetchUsers = useCallback(() => {
    api.get<PaginatedResponse<User>>('/users', { params: { page, per_page: USERS_PER_PAGE } })
      .then(({ data }) => {
        setUsers(data.data);
        setTotalPages(data.total_pages);
      })
      .catch(() => showSnackbar('Error al cargar usuarios', 'error'));
  }, [page, showSnackbar]);

  useEffect(() => {
    fetchUsers();
  }, [fetchUsers]);

  useEffect(() => {
    setLoading(true);
    api.get<AdminStats>('/admin/stats')
      .then(({ data }) => setStats(data))
      .catch(() => showSnackbar('Error al cargar estadísticas', 'error'))
      .finally(() => setLoading(false));
  }, [showSnackbar]);

  const promoteToVerifier = async (userId: string) => {
    try {
      await api.patch(`/users/${userId}/role`, { role: 'verifier' });
      setUsers((prev) => prev.map((u) => u.id === userId ? { ...u, role: 'verifier' } : u));
      showSnackbar('Usuario promovido a verificador', 'success');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error', 'error');
    }
  };

  const suspendUser = async (userId: string) => {
    setSuspendTarget(null);
    try {
      await api.patch(`/users/${userId}/status`, { status: 'banned' });
      setUsers((prev) => prev.map((u) => u.id === userId ? { ...u, status: 'banned' } : u));
      showSnackbar('Usuario suspendido', 'info');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error', 'error');
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Panel de Administración</Typography>

      <Grid container spacing={2} sx={{ mb: 4 }}>
        <Grid size={{ xs: 6, md: 3 }}>
          <StatCard icon={<PeopleIcon sx={{ fontSize: 32 }} />} label="Usuarios" value={stats.total_users} />
        </Grid>
        <Grid size={{ xs: 6, md: 3 }}>
          <StatCard icon={<CasinoIcon sx={{ fontSize: 32 }} />} label="Rifas" value={stats.total_raffles} />
        </Grid>
        <Grid size={{ xs: 6, md: 3 }}>
          <StatCard icon={<ConfirmationNumberIcon sx={{ fontSize: 32 }} />} label="Boletos" value={stats.total_tickets} />
        </Grid>
        <Grid size={{ xs: 6, md: 3 }}>
          <StatCard icon={<GavelIcon sx={{ fontSize: 32 }} />} label="Pendientes" value={stats.pending_tickets} />
        </Grid>
      </Grid>

      <Typography variant="h6" gutterBottom>Usuarios</Typography>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Nombre</TableCell>
            <TableCell>Email</TableCell>
            <TableCell>Rol</TableCell>
            <TableCell>Estado</TableCell>
            <TableCell>Acciones</TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          {users.map((user) => (
            <TableRow key={user.id}>
              <TableCell>{user.name}</TableCell>
              <TableCell>{user.email}</TableCell>
              <TableCell><Chip label={roleLabel(user.role)} size="small" /></TableCell>
              <TableCell>{user.status}</TableCell>
              <TableCell sx={{ display: 'flex', gap: 1 }}>
                {user.role !== 'verifier' && user.role !== 'admin' && (
                  <Button size="small" onClick={() => promoteToVerifier(user.id)}>Hacer Verificador</Button>
                )}
                {user.status === 'active' && (
                  <Button size="small" color="error" onClick={() => setSuspendTarget(user)}>Suspender</Button>
                )}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
      {totalPages > 1 && (
        <Pagination count={totalPages} page={page} onChange={(_, p) => setPage(p)} sx={{ mt: 2, display: 'flex', justifyContent: 'center' }} />
      )}

      <Dialog open={!!suspendTarget} onClose={() => setSuspendTarget(null)}>
        <DialogTitle>¿Suspender usuario?</DialogTitle>
        <DialogContent>
          <Typography>
            Esto impedirá que <strong>{suspendTarget?.name}</strong> acceda a la plataforma.
          </Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setSuspendTarget(null)}>Cancelar</Button>
          <Button color="error" variant="contained" onClick={() => suspendTarget && suspendUser(suspendTarget.id)}>
            Suspender
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
