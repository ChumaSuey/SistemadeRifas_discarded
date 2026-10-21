import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Typography, Button, Box, TextField, Chip, Alert, FormControl, InputLabel, Select, MenuItem, CircularProgress } from '@mui/material';
import SettingsIcon from '@mui/icons-material/Settings';
import api from '../api/client';
import type { RaffleDetail, Ticket } from '../types';
import { useAuth } from '../context/AuthContext';
import { useSnackbar } from '../context/SnackbarContext';
import { formatCurrency, statusLabel, statusColor } from '../utils/formatters';

export default function RaffleDetailPage() {
  const { id } = useParams<{ id: string }>();
  const { user, isAuthenticated } = useAuth();
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();
  const [raffle, setRaffle] = useState<RaffleDetail | null>(null);
  const [quantity, setQuantity] = useState(1);
  const [paymentMethod, setPaymentMethod] = useState('manual_transfer');
  const [buying, setBuying] = useState(false);
  const [boughtTickets, setBoughtTickets] = useState<Ticket[] | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    setLoading(true);
    api.get<RaffleDetail>(`/raffles/${id}`)
      .then(({ data }) => setRaffle(data))
      .catch(() => showSnackbar('Error al cargar la rifa', 'error'))
      .finally(() => setLoading(false));
  }, [id, showSnackbar]);

  const handleBuy = async () => {
    setBuying(true);
    try {
      const { data } = await api.post<Ticket[]>(`/raffles/${id}/tickets`, { quantity, payment_method: paymentMethod });
      setBoughtTickets(data);
      const { data: updated } = await api.get<RaffleDetail>(`/raffles/${id}`);
      setRaffle(updated);
      showSnackbar(`Compraste ${data.length} boleto(s)`, 'success');
    } catch (err: unknown) {
      const msg = (err as { response?: { data?: { detail?: string; message?: string } } })?.response?.data?.detail
        || (err as { response?: { data?: { message?: string } } })?.response?.data?.message
        || 'Error al comprar tickets';
      showSnackbar(msg, 'error');
    } finally {
      setBuying(false);
    }
  };

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!raffle) return <Typography>Cargando...</Typography>;

  const available = raffle.tickets_total - raffle.tickets_sold;
  const isOwner = user && raffle.raffler_id === user.id;

  return (
    <Box>
      {isOwner && (
        <Button
          variant="outlined"
          startIcon={<SettingsIcon />}
          onClick={() => navigate(`/raffles/${raffle.id}/manage`)}
          sx={{ mb: 2 }}
        >
          Gestionar Rifa
        </Button>
      )}

      <Typography variant="h4">{raffle.title}</Typography>
      <Box sx={{ display: 'flex', gap: 2, my: 1, flexWrap: 'wrap' }}>
        <Chip label={statusLabel(raffle.status)} color={statusColor(raffle.status)} />
        <Chip label={formatCurrency(raffle.ticket_price, raffle.currency)} variant="outlined" />
        <Chip label={`${raffle.tickets_sold}/${raffle.tickets_total} boletos`} variant="outlined" />
      </Box>
      <Typography sx={{ mt: 2, whiteSpace: 'pre-wrap' }}>{raffle.description}</Typography>

      {raffle.status === 'active' && isAuthenticated && available > 0 && !isOwner && (
        <Box sx={{ mt: 3, p: 2, border: '1px solid', borderColor: 'divider', borderRadius: 1 }}>
          <Typography variant="h6" gutterBottom>Comprar Boletos</Typography>
          {boughtTickets ? (
            <>
              <Alert severity="success">
                Compraste {boughtTickets.length} boleto(s): {boughtTickets.map((t) => `#${t.ticket_number}`).join(', ')}
              </Alert>
              {paymentMethod === 'manual_transfer' && (
                <Alert severity="info" sx={{ mt: 1 }}>
                  Sube tu comprobante de pago en la sección <strong>Mis Tickets</strong>
                </Alert>
              )}
            </>
          ) : (
            <Box sx={{ display: 'flex', gap: 2, alignItems: 'flex-end', flexWrap: 'wrap' }}>
              <TextField type="number" label="Cantidad" value={quantity} onChange={(e) => setQuantity(Math.max(1, Math.min(+e.target.value, available)))} size="small" sx={{ width: 100 }} />
              <FormControl size="small" sx={{ minWidth: 160 }}>
                <InputLabel>Método de pago</InputLabel>
                <Select value={paymentMethod} label="Método de pago" onChange={(e) => setPaymentMethod(e.target.value)}>
                  <MenuItem value="manual_transfer">Transferencia</MenuItem>
                  <MenuItem value="stripe">Stripe (Pendiente)</MenuItem>
                </Select>
              </FormControl>
              <Button variant="contained" onClick={handleBuy} disabled={buying}>
                {buying ? 'Comprando...' : 'Comprar'}
              </Button>
            </Box>
          )}
        </Box>
      )}

      {raffle.status === 'finished' && raffle.winner_id && (
        <Alert severity="success" sx={{ mt: 3 }}>
          Ganador del sorteo: Ticket #{raffle.tickets?.find((t) => t.id === raffle.winning_ticket_id)?.ticket_number || 'N/A'}
        </Alert>
      )}
    </Box>
  );
}
