import { useState, useEffect } from 'react';
import { Typography, Box, Button, CircularProgress } from '@mui/material';
import UploadIcon from '@mui/icons-material/Upload';
import ConfirmationNumberIcon from '@mui/icons-material/ConfirmationNumber';
import api from '../api/client';
import type { Ticket, PaginatedResponse } from '../types';
import { useSnackbar } from '../context/SnackbarContext';
import TicketCard from '../components/TicketCard';

export default function MyTicketsPage() {
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [loading, setLoading] = useState(true);
  const { showSnackbar } = useSnackbar();

  useEffect(() => {
    setLoading(true);
    api.get<PaginatedResponse<Ticket>>('/tickets/my')
      .then(({ data }) => setTickets(data.data))
      .catch(() => showSnackbar('Error al cargar tus boletos', 'error'))
      .finally(() => setLoading(false));
  }, [showSnackbar]);

  const handleUploadProof = async (ticketId: string) => {
    const url = prompt('URL del comprobante:');
    if (!url) return;
    try {
      await api.post(`/tickets/${ticketId}/upload-proof`, { payment_proof_url: url });
      setTickets((prev) => prev.map((t) => t.id === ticketId ? { ...t, payment_proof_url: url } : t));
      showSnackbar('Comprobante subido', 'success');
    } catch {
      showSnackbar('Error al subir comprobante', 'error');
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
      <Typography variant="h4" gutterBottom>Mis Boletos</Typography>

      {tickets.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 8 }}>
          <ConfirmationNumberIcon sx={{ fontSize: 64, color: 'text.disabled', mb: 2 }} />
          <Typography variant="h6" color="text.secondary" gutterBottom>
            No has comprado boletos aún
          </Typography>
          <Typography color="text.secondary">
            Explora rifas activas y compra tus boletos para participar
          </Typography>
        </Box>
      ) : (
        tickets.map((ticket) => (
          <Box key={ticket.id} sx={{ mb: 1 }}>
            <TicketCard ticket={ticket} />
            {ticket.status === 'pending' && ticket.payment_method === 'manual_transfer' && !ticket.payment_proof_url && (
              <Button size="small" startIcon={<UploadIcon />} onClick={() => handleUploadProof(ticket.id)} sx={{ ml: 2 }}>
                Subir Comprobante
              </Button>
            )}
          </Box>
        ))
      )}
    </Box>
  );
}
