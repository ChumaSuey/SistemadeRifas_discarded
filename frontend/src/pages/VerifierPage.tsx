import { useState, useEffect } from 'react';
import { Typography, Box, Button, CircularProgress, Pagination } from '@mui/material';
import CheckIcon from '@mui/icons-material/Check';
import CloseIcon from '@mui/icons-material/Close';
import GavelIcon from '@mui/icons-material/Gavel';
import api from '../api/client';
import type { Ticket } from '../types';
import { useSnackbar } from '../context/SnackbarContext';
import TicketCard from '../components/TicketCard';

const PER_PAGE = 10;

export default function VerifierPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const [loadingAction, setLoadingAction] = useState(false);
  const [page, setPage] = useState(1);
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    setLoading(true);
    api.get<{ data: Ticket[] }>('/tickets/pending')
      .then(({ data }) => setTickets(data.data || data as unknown as Ticket[] || []))
      .catch(() => {
        showSnackbar('Error al cargar tickets pendientes', 'error');
        setTickets([]);
      })
      .finally(() => setLoading(false));
  }, [showSnackbar]);

  const handleVerify = async (ticketId: string) => {
    setLoadingAction(true);
    try {
      await api.post(`/tickets/${ticketId}/verify`);
      setTickets((prev) => prev.filter((t) => t.id !== ticketId));
      showSnackbar('Ticket aprobado', 'success');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al aprobar', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const handleReject = async (ticketId: string, reason: string) => {
    setLoadingAction(true);
    try {
      await api.post(`/tickets/${ticketId}/reject`, { reason });
      setTickets((prev) => prev.filter((t) => t.id !== ticketId));
      showSnackbar('Ticket rechazado', 'info');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al rechazar', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const totalPages = Math.ceil(tickets.length / PER_PAGE);
  const paginatedTickets = tickets.slice((page - 1) * PER_PAGE, page * PER_PAGE);

  if (loading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  return (
    <Box>
      <Typography variant="h4" gutterBottom>Panel de Verificación</Typography>

      {tickets.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <GavelIcon sx={{ fontSize: 64, color: 'text.disabled', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No hay boletos pendientes
          </Typography>
          <Typography color="text.secondary">
            Los boletos pendientes de verificación aparecerán aquí
          </Typography>
        </Box>
      ) : (
        <>
          {paginatedTickets.map((ticket) => (
            <Box key={ticket.id} sx={{ p: 2, mb: 1, border: '1px solid', borderColor: 'divider', borderRadius: 1, display: 'flex', alignItems: 'center', gap: 2, flexWrap: 'wrap' }}>
              <Box sx={{ flex: 1, minWidth: 200 }}>
                <TicketCard ticket={ticket} />
              </Box>
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button variant="contained" color="success" size="small" startIcon={<CheckIcon />}
                  onClick={() => handleVerify(ticket.id)} disabled={loadingAction}>
                  Aprobar
                </Button>
                <Button variant="outlined" color="error" size="small" startIcon={<CloseIcon />}
                  onClick={() => {
                    const reason = prompt('Razón de rechazo:');
                    if (reason) handleReject(ticket.id, reason);
                  }}
                  disabled={loadingAction}>
                  Rechazar
                </Button>
              </Box>
            </Box>
          ))}
          {totalPages > 1 && (
            <Pagination count={totalPages} page={page} onChange={(_, p) => setPage(p)} sx={{ mt: 2, display: 'flex', justifyContent: 'center' }} />
          )}
        </>
      )}
    </Box>
  );
}
