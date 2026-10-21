import { useState, useEffect, useRef, useCallback } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Typography, Button, Box, Chip, Alert, CircularProgress, IconButton } from '@mui/material';
import ContentCopyIcon from '@mui/icons-material/ContentCopy';
import api from '../api/client';
import type { SortRevealResponse, RaffleDetail } from '../types';
import { useSnackbar } from '../context/SnackbarContext';
import { formatCurrency, statusLabel, statusColor } from '../utils/formatters';

function shuffleArray(arr: number[]) {
  const a = [...arr];
  for (let i = a.length - 1; i > 0; i--) {
    const j = Math.floor(Math.random() * (i + 1));
    [a[i], a[j]] = [a[j], a[i]];
  }
  return a;
}

export default function LiveDrawPage() {
  const { id } = useParams<{ id: string }>();
  const navigate = useNavigate();
  const { showSnackbar } = useSnackbar();
  const [raffle, setRaffle] = useState<RaffleDetail | null>(null);
  const [result, setResult] = useState<SortRevealResponse | null>(null);
  const [loading, setLoading] = useState(false);
  const [animating, setAnimating] = useState(false);
  const [displayNumber, setDisplayNumber] = useState<number | null>(null);
  const [pageLoading, setPageLoading] = useState(true);
  const intervalRef = useRef<ReturnType<typeof setInterval> | null>(null);

  useEffect(() => {
    api.get<RaffleDetail>(`/raffles/${id}`)
      .then(({ data }) => setRaffle(data))
      .catch(() => showSnackbar('Error al cargar la rifa', 'error'))
      .finally(() => setPageLoading(false));
  }, [id, showSnackbar]);

  const runAnimation = useCallback((finalNumber: number, ticketNumbers: number[]) => {
    setAnimating(true);
    setDisplayNumber(0);
    const shuffled = shuffleArray(ticketNumbers);
    let idx = 0;
    const fastPhase = 8;
    let tick = 0;

    intervalRef.current = setInterval(() => {
      tick++;
      if (tick <= fastPhase) {
        idx = (idx + 1) % shuffled.length;
        setDisplayNumber(shuffled[idx]);
      } else if (tick <= fastPhase + 5) {
        const progress = (tick - fastPhase) / 5;
        if (Math.random() < 1 - progress) {
          idx = (idx + 1) % shuffled.length;
          setDisplayNumber(shuffled[idx]);
        }
      } else if (tick === fastPhase + 6) {
        setDisplayNumber(finalNumber);
        setAnimating(false);
        if (intervalRef.current) clearInterval(intervalRef.current);
      }
    }, 100);
  }, []);

  const handleSort = async () => {
    setLoading(true);
    try {
      const { data: commitData } = await api.post<{ commit_hash: string; tickets_count: number }>(`/raffles/${id}/commit`);

      const { data: revealData } = await api.post<SortRevealResponse>(`/raffles/${id}/recompute-and-reveal`, commitData);

      setResult(revealData);
      const verified = raffle?.tickets?.filter((t) => t.status === 'verified') || [];
      const nums = verified.length > 0
        ? verified.map((t) => t.ticket_number)
        : Array.from({ length: raffle?.tickets_sold || 50 }, (_, i) => i + 1);
      runAnimation(revealData.ticket_number, nums);
      showSnackbar('Sorteo ejecutado', 'success');
    } catch {
      showSnackbar('Error al ejecutar sorteo', 'error');
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    return () => { if (intervalRef.current) clearInterval(intervalRef.current); };
  }, []);

  if (pageLoading) {
    return (
      <Box sx={{ display: 'flex', justifyContent: 'center', py: 8 }}>
        <CircularProgress />
      </Box>
    );
  }

  if (!raffle) return <Typography>Rifa no encontrada.</Typography>;

  const verifiedCount = raffle.tickets?.filter((t) => t.status === 'verified').length || 0;

  return (
    <Box sx={{ textAlign: 'center', py: 4 }}>
      <Typography variant="h3" gutterBottom>{raffle.title}</Typography>
      <Box sx={{ display: 'flex', gap: 1, justifyContent: 'center', mb: 3, flexWrap: 'wrap' }}>
        <Chip label={statusLabel(raffle.status)} color={statusColor(raffle.status)} />
        <Chip label={`${verifiedCount} boletos verificados`} variant="outlined" />
        <Chip label={formatCurrency(raffle.ticket_price, raffle.currency)} variant="outlined" />
      </Box>

      {animating || (displayNumber !== null && animating) ? (
        <Box sx={{ mt: 6, p: 4, border: '2px solid', borderColor: 'warning.main', borderRadius: 2 }}>
          <Typography variant="h5" color="text.secondary" gutterBottom>
            Seleccionando ganador...
          </Typography>
          <Typography variant="h1" sx={{ fontSize: { xs: '4rem', md: '8rem' }, fontWeight: 700, color: 'primary.main', transition: 'all .05s ease' }}>
            #{displayNumber ?? '?'}
          </Typography>
        </Box>
      ) : null}

      {result && !animating ? (
        <Box sx={{ mt: 4, p: 4, border: '3px solid', borderColor: 'success.main', borderRadius: 2, background: 'linear-gradient(135deg, #e8f5e9, #c8e6c9)' }}>
          <Typography variant="h3" color="success.main" gutterBottom sx={{ fontWeight: 700 }}>
            ¡Ganador!
          </Typography>
          <Typography variant="h1" sx={{ fontSize: { xs: '4rem', md: '8rem' }, fontWeight: 800, color: '#1b5e20' }}>
            #{result.ticket_number}
          </Typography>
          <Box sx={{ mt: 2 }}>
            <Typography color="text.secondary" component="span" sx={{ wordBreak: 'break-all' }}>
              Hash público: {result.commit_hash}
            </Typography>
            <IconButton
              size="small"
              onClick={() => { navigator.clipboard.writeText(result.commit_hash); showSnackbar('Hash copiado', 'info'); }}
              sx={{ ml: 0.5 }}
            >
              <ContentCopyIcon fontSize="small" />
            </IconButton>
          </Box>
          <Box sx={{ mt: 0.5 }}>
            <Typography variant="caption" component="span">
              Semilla: {result.seed}
            </Typography>
            <IconButton
              size="small"
              onClick={() => { navigator.clipboard.writeText(result.seed); showSnackbar('Semilla copiada', 'info'); }}
              sx={{ ml: 0.5 }}
            >
              <ContentCopyIcon fontSize="small" />
            </IconButton>
          </Box>
          <Typography color="text.secondary" sx={{ mt: 1 }}>
            Fecha del sorteo: {result.draw_date}
          </Typography>
          <Box sx={{ mt: 3, display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
            <Button variant="outlined" onClick={() => navigate(`/raffles/${id}/manage`)}>Volver a gestión</Button>
            <Button variant="contained" onClick={() => navigate(`/raffles/${id}/manage`)}>Re-Sortear</Button>
          </Box>
        </Box>
      ) : !animating ? (
        loading ? (
          <Box sx={{ mt: 8 }}>
            <CircularProgress size={80} />
            <Typography sx={{ mt: 2 }}>Ejecutando sorteo...</Typography>
          </Box>
        ) : raffle.status === 'finished' ? (
          <Alert severity="info" sx={{ mt: 4 }}>Esta rifa ya finalizó. Puedes ver los resultados en la página de gestión.</Alert>
        ) : (
          <Box sx={{ mt: 8 }}>
            <Typography variant="h6" color="text.secondary" gutterBottom>
              {verifiedCount} boletos verificados listos para el sorteo
            </Typography>
            <Button variant="contained" size="large" onClick={handleSort} disabled={loading} sx={{ mt: 2, px: 6, py: 1.5 }}>
              {loading ? 'Procesando...' : 'Ejecutar Sorteo'}
            </Button>
          </Box>
        )
      ) : null}
    </Box>
  );
}
