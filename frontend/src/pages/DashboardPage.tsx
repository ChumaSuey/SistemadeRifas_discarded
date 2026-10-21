import { useState, useEffect } from 'react';
import { Typography, Button, Box, TextField, Dialog, DialogTitle, DialogContent, DialogActions, Grid, FormControl, InputLabel, Select, MenuItem, CircularProgress } from '@mui/material';
import AddIcon from '@mui/icons-material/Add';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';
import type { Raffle, PaginatedResponse } from '../types';
import RaffleCard from '../components/RaffleCard';
import { useSnackbar } from '../context/SnackbarContext';

const defaultForm = { title: '', description: '', ticket_price: 5, currency: 'USD', tickets_total: 100, payment_methods: ['manual_transfer'] };

export default function DashboardPage() {
  const [raffles, setRaffles] = useState<Raffle[]>([]);
  const [open, setOpen] = useState(false);
  const [form, setForm] = useState(defaultForm);
  const [creating, setCreating] = useState(false);
  const [loading, setLoading] = useState(true);
  const [refreshKey, setRefreshKey] = useState(0);
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    setLoading(true);
    api.get<PaginatedResponse<Raffle>>('/raffles/my')
      .then(({ data }) => setRaffles(data.data))
      .catch(() => showSnackbar('Error al cargar tus rifas', 'error'))
      .finally(() => setLoading(false));
  }, [showSnackbar, refreshKey]);

  const handleCreate = async () => {
    setCreating(true);
    try {
      const { data } = await api.post<Raffle>('/raffles', form);
      setOpen(false);
      setForm(defaultForm);
      setRefreshKey((k) => k + 1);
      showSnackbar('Rifa creada exitosamente', 'success');
      navigate(`/raffles/${data.id}/manage`);
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { detail?: string; message?: string } } })?.response?.data?.detail
        || (err as { response?: { data?: { message?: string } } })?.response?.data?.message
        || 'Error al crear rifa';
      showSnackbar(msg, 'error');
    } finally {
      setCreating(false);
    }
  };

  return (
    <Box>
      <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 3 }}>
        <Typography variant="h4">Mis Rifas</Typography>
        <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)}>Nueva Rifa</Button>
      </Box>

      {loading ? (
        <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
          <CircularProgress />
        </Box>
      ) : (
        <>
          <Grid container spacing={2}>
            {raffles.map((r) => (
              <Grid size={{ xs: 12, sm: 6, md: 4 }} key={r.id}>
                <RaffleCard raffle={r} />
              </Grid>
            ))}
          </Grid>
          {raffles.length === 0 && (
            <Box sx={{ textAlign: 'center', py: 8 }}>
              <Typography variant="h6" color="text.secondary" gutterBottom>
                No has creado rifas aún
              </Typography>
              <Button variant="contained" startIcon={<AddIcon />} onClick={() => setOpen(true)} sx={{ mt: 2 }}>
                Crear mi primera rifa
              </Button>
            </Box>
          )}
        </>
      )}

      <Dialog open={open} onClose={() => setOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Nueva Rifa</DialogTitle>
        <DialogContent>
          <TextField fullWidth label="Título" value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} margin="normal" />
          <TextField fullWidth label="Descripción" value={form.description} onChange={(e) => setForm({ ...form, description: e.target.value })} margin="normal" multiline rows={3} />
          <TextField fullWidth label="Precio por boleto" type="number" value={form.ticket_price} onChange={(e) => setForm({ ...form, ticket_price: +e.target.value })} margin="normal" />
          <FormControl fullWidth margin="normal">
            <InputLabel>Moneda</InputLabel>
            <Select value={form.currency} label="Moneda" onChange={(e) => setForm({ ...form, currency: e.target.value })}>
              <MenuItem value="USD">USD</MenuItem>
              <MenuItem value="VES">VES</MenuItem>
              <MenuItem value="EUR">EUR</MenuItem>
            </Select>
          </FormControl>
          <TextField fullWidth label="Total de boletos" type="number" value={form.tickets_total} onChange={(e) => setForm({ ...form, tickets_total: +e.target.value })} margin="normal" />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setOpen(false)}>Cancelar</Button>
          <Button variant="contained" onClick={handleCreate} disabled={creating || !form.title}>
            {creating ? 'Creando...' : 'Crear'}
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
