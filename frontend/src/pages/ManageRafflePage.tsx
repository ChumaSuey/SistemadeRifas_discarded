import { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Typography, Button, Box, Chip, Alert, CircularProgress, TextField, Dialog, DialogTitle, DialogContent, DialogActions, Pagination } from '@mui/material';
import PlayArrowIcon from '@mui/icons-material/PlayArrow';
import CancelIcon from '@mui/icons-material/Cancel';
import PublishIcon from '@mui/icons-material/Publish';
import PanToolIcon from '@mui/icons-material/PanTool';
import api from '../api/client';
import type { RaffleDetail, SortRevealResponse, Ticket, SortCommitResponse } from '../types';
import { useAuth } from '../context/AuthContext';
import { useSnackbar } from '../context/SnackbarContext';
import { formatCurrency, statusLabel, statusColor } from '../utils/formatters';
import TicketCard from '../components/TicketCard';

export default function RaffleManagePage() {
  const { id } = useParams<{ id: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();
  const [raffle, setRaffle] = useState<RaffleDetail | null>(null);
  const [tickets, setTickets] = useState<Ticket[]>([]);
  const [sortResult, setSortResult] = useState<SortRevealResponse | null>(null);
  const [commitPhase, setCommitPhase] = useState<SortCommitResponse | null>(null);
  const [resortReason, setResortReason] = useState('');
  const [loadingAction, setLoadingAction] = useState(false);
  const [pageLoading, setPageLoading] = useState(true);
  const [confirmCancelOpen, setConfirmCancelOpen] = useState(false);
  const [ticketsPage, setTicketsPage] = useState(1);
  const TICKETS_PER_PAGE = 10;

  const owner = raffle && user && raffle.raffler_id === user.id;

  useEffect(() => {
    setPageLoading(true);
    Promise.all([
      api.get<RaffleDetail>(`/raffles/${id}`).then(({ data }) => setRaffle(data)),
      api.get<{ data: Ticket[] }>(`/raffles/${id}/tickets`).then(({ data }) => setTickets(data.data || data)),
    ])
      .catch(() => showSnackbar('Error al cargar la rifa', 'error'))
      .finally(() => setPageLoading(false));
  }, [id, showSnackbar]);

  const handlePublish = async () => {
    setLoadingAction(true);
    try {
      await api.post(`/raffles/${id}/publish`);
      setRaffle((prev) => prev ? { ...prev, status: 'active' } : null);
      showSnackbar('Rifa publicada exitosamente', 'success');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al publicar', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const handleCancel = async () => {
    setConfirmCancelOpen(false);
    setLoadingAction(true);
    try {
      await api.post(`/raffles/${id}/cancel`);
      setRaffle((prev) => prev ? { ...prev, status: 'cancelled' } : null);
      showSnackbar('Rifa cancelada', 'info');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al cancelar', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const handleCommit = async () => {
    setLoadingAction(true);
    try {
      const { data } = await api.post<SortCommitResponse>(`/raffles/${id}/commit`);
      setCommitPhase(data);
      showSnackbar('Hash de compromiso generado', 'info');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al generar commit', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const handleSort = async () => {
    setLoadingAction(true);
    try {
      const { data } = await api.post<SortRevealResponse>(`/raffles/${id}/sort`);
      setSortResult(data);
      setRaffle((prev) => prev ? { ...prev, status: 'finished', winner_id: data.winner_user_id } : null);
      showSnackbar('Sorteo ejecutado', 'success');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al ejecutar sorteo', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  const handleResort = async () => {
    setLoadingAction(true);
    try {
      const { data } = await api.post<SortRevealResponse>(`/raffles/${id}/resort`, { disqualify_reason: resortReason || 'Re-sorteo solicitado' });
      setSortResult(data);
      setResortReason('');
      showSnackbar('Re-sorteo ejecutado', 'info');
    } catch (err: unknown) {
      showSnackbar((err as { response?: { data?: { detail?: string } } })?.response?.data?.detail || 'Error al re-sortear', 'error');
    } finally {
      setLoadingAction(false);
    }
  };

  if (pageLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!raffle) return <Typography>Rifa no encontrada.</Typography>;
  if (!owner) return <Typography>No tienes permiso para gestionar esta rifa.</Typography>;

  const verifiedTickets = tickets.filter((t) => t.status === 'verified');

  return (
    <Box>
      <Typography variant="h4">{raffle.title}</Typography>
      <Box sx={{ display: 'flex', gap: 1, mt: 1, mb: 2, flexWrap: 'wrap' }}>
        <Chip label={statusLabel(raffle.status)} color={statusColor(raffle.status)} />
        <Chip label={`${raffle.tickets_sold}/${raffle.tickets_total} boletos`} variant="outlined" />
        <Chip label={formatCurrency(raffle.ticket_price, raffle.currency)} variant="outlined" />
      </Box>

      {sortResult && (
        <Alert severity="success" sx={{ mb: 2 }}>
          Ganador: Ticket #{sortResult.ticket_number} — Hash: {sortResult.commit_hash.slice(0, 16)}...
        </Alert>
      )}

      {commitPhase && (
        <Alert severity="info" sx={{ mb: 2 }}>
          Commit generado — Hash: {commitPhase.commit_hash.slice(0, 24)}... — {commitPhase.tickets_count} boletos verificados — {new Date(commitPhase.timestamp * 1000).toLocaleString()}
        </Alert>
      )}

      <Box sx={{ display: 'flex', gap: 1, mb: 3, flexWrap: 'wrap' }}>
        {raffle.status === 'draft' && (
          <Button variant="contained" color="success" startIcon={<PublishIcon />} onClick={handlePublish} disabled={loadingAction}>
            Publicar Rifa
          </Button>
        )}
        {raffle.status === 'active' && (
          <>
            {!commitPhase && (
              <Button variant="outlined" onClick={handleCommit} disabled={loadingAction}>
                Generar Commit
              </Button>
            )}
            <Button variant="contained" startIcon={<PlayArrowIcon />} onClick={handleSort} disabled={loadingAction}>
              Ejecutar Sorteo
            </Button>
            <Button variant="outlined" color="error" startIcon={<CancelIcon />} onClick={() => setConfirmCancelOpen(true)} disabled={loadingAction}>
              Cancelar Rifa
            </Button>
          </>
        )}
        {raffle.status === 'finished' && (
          <Button variant="outlined" startIcon={<PanToolIcon />} onClick={handleResort} disabled={loadingAction}>
            Re-Sortear
          </Button>
        )}
        <Button onClick={() => navigate(`/raffles/${id}/live`)} variant="outlined">
          Vista en Vivo
        </Button>
      </Box>

      {raffle.status === 'finished' && (
        <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
          <TextField value={resortReason} onChange={(e) => setResortReason(e.target.value)} placeholder="Razón de re-sorteo" sx={{ flex: 1, minWidth: 200 }} size="small" />
          <Button variant="contained" color="warning" startIcon={<PanToolIcon />} onClick={handleResort} disabled={loadingAction}>
            Re-Sortear
          </Button>
        </Box>
      )}

      <Typography variant="h6" gutterBottom>Boletos Verificados ({verifiedTickets.length})</Typography>
      {verifiedTickets.length === 0 ? (
        <Box sx={{ textAlign: 'center', py: 4 }}>
          <Typography color="text.secondary">Aún no hay boletos verificados.</Typography>
        </Box>
      ) : (
        <>
          {verifiedTickets.slice((ticketsPage - 1) * TICKETS_PER_PAGE, ticketsPage * TICKETS_PER_PAGE).map((t) => (
            <TicketCard key={t.id} ticket={t} />
          ))}
          {verifiedTickets.length > TICKETS_PER_PAGE && (
            <Pagination
              count={Math.ceil(verifiedTickets.length / TICKETS_PER_PAGE)}
              page={ticketsPage}
              onChange={(_, p) => setTicketsPage(p)}
              sx={{ mt: 2, display: 'flex', justifyContent: 'center' }}
            />
          )}
        </>
      )}

      <Dialog open={confirmCancelOpen} onClose={() => setConfirmCancelOpen(false)}>
        <DialogTitle>¿Cancelar rifa?</DialogTitle>
        <DialogContent>
          <Typography>Esta acción no se puede deshacer. Todos los boletos comprados serán cancelados.</Typography>
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setConfirmCancelOpen(false)}>Volver</Button>
          <Button color="error" variant="contained" onClick={handleCancel}>Cancelar Rifa</Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}
