import { useState, useEffect } from 'react';
import { Typography, Grid, Box, Pagination, CircularProgress, Button, Container, Paper } from '@mui/material';
import CasinoIcon from '@mui/icons-material/Casino';
import PaymentIcon from '@mui/icons-material/Payment';
import EmojiEventsIcon from '@mui/icons-material/EmojiEvents';
import ArrowForwardIcon from '@mui/icons-material/ArrowForward';
import { useNavigate } from 'react-router-dom';
import api from '../api/client';
import type { Raffle, PaginatedResponse } from '../types';
import RaffleCard from '../components/RaffleCard';
import { useAuth } from '../context/AuthContext';
import { useSnackbar } from '../context/SnackbarContext';

const steps = [
  { icon: <CasinoIcon sx={{ fontSize: 48 }} />, title: 'Encuentra una rifa', desc: 'Explora rifas activas de streamers y creadores de contenido. Elige la que más te guste.' },
  { icon: <PaymentIcon sx={{ fontSize: 48 }} />, title: 'Compra tus boletos', desc: 'Selecciona la cantidad de boletos y realiza el pago por transferencia o tarjeta.' },
  { icon: <EmojiEventsIcon sx={{ fontSize: 48 }} />, title: 'Gana el premio', desc: 'El sorteo es transparente y verificable. Si ganas, recibirás el premio automáticamente.' },
];

export default function HomePage() {
  const [raffles, setRaffles] = useState<Raffle[]>([]);
  const [page, setPage] = useState(1);
  const [totalPages, setTotalPages] = useState(1);
  const [loading, setLoading] = useState(true);
  const { showSnackbar } = useSnackbar();
  const { isAuthenticated } = useAuth();
  const navigate = useNavigate();

  useEffect(() => {
    setLoading(true);
    api.get<PaginatedResponse<Raffle>>('/raffles', { params: { page, per_page: 12 } })
      .then(({ data }) => {
        setRaffles(data.data);
        setTotalPages(data.total_pages);
      })
      .catch(() => showSnackbar('Error al cargar rifas', 'error'))
      .finally(() => setLoading(false));
  }, [page, showSnackbar]);

  return (
    <Box>
      <Box
        sx={{
          background: 'linear-gradient(135deg, #1565c0 0%, #1976d2 40%, #42a5f5 100%)',
          color: 'white',
          py: { xs: 6, md: 10 },
          px: 2,
          textAlign: 'center',
          borderRadius: { xs: 0, md: 3 },
          mt: { xs: -3, md: 0 },
          mx: { md: -2 },
        }}
      >
        <Container maxWidth="md">
          <Typography variant="h3" sx={{ fontWeight: 800, mb: 2, fontSize: { xs: '2rem', md: '3rem' } }}>
            Rifas transparentes para creadores
          </Typography>
          <Typography variant="h6" sx={{ mb: 4, opacity: 0.9, fontWeight: 400, fontSize: { xs: '1rem', md: '1.25rem' } }}>
            Crea, gestiona y participa en rifas con sorteos verificables mediante hash criptográfico. Transparencia garantizada.
          </Typography>
          <Box sx={{ display: 'flex', gap: 2, justifyContent: 'center', flexWrap: 'wrap' }}>
            {isAuthenticated ? (
              <Button
                variant="contained"
                size="large"
                onClick={() => navigate('/dashboard')}
                sx={{ bgcolor: 'white', color: '#1565c0', fontWeight: 700, px: 4, py: 1.5, '&:hover': { bgcolor: '#e3f2fd' } }}
                endIcon={<ArrowForwardIcon />}
              >
                Ir al Dashboard
              </Button>
            ) : (
              <>
                <Button
                  variant="contained"
                  size="large"
                  onClick={() => navigate('/login')}
                  sx={{ bgcolor: 'white', color: '#1565c0', fontWeight: 700, px: 4, py: 1.5, '&:hover': { bgcolor: '#e3f2fd' } }}
                  endIcon={<ArrowForwardIcon />}
                >
                  Comenzar ahora
                </Button>
                <Button
                  variant="outlined"
                  size="large"
                  onClick={() => document.getElementById('raffles-section')?.scrollIntoView({ behavior: 'smooth' })}
                  sx={{ borderColor: 'white', color: 'white', fontWeight: 700, px: 4, py: 1.5, '&:hover': { borderColor: 'white', bgcolor: 'rgba(255,255,255,0.1)' } }}
                >
                  Ver rifas
                </Button>
              </>
            )}
          </Box>
        </Container>
      </Box>

      <Container maxWidth="md" sx={{ py: { xs: 6, md: 10 } }}>
        <Typography variant="h4" align="center" sx={{ fontWeight: 700, mb: 6 }}>
          Cómo funciona
        </Typography>
        <Grid container spacing={4}>
          {steps.map((step, idx) => (
            <Grid size={{ xs: 12, md: 4 }} key={idx}>
              <Paper
                variant="outlined"
                sx={{
                  p: 4,
                  textAlign: 'center',
                  height: '100%',
                  borderColor: 'primary.main',
                  borderWidth: 1.5,
                  transition: 'transform 0.2s, box-shadow 0.2s',
                  '&:hover': { transform: 'translateY(-4px)', boxShadow: 4 },
                }}
              >
                <Box sx={{ color: 'primary.main', mb: 2 }}>{step.icon}</Box>
                <Typography variant="h6" sx={{ fontWeight: 700, mb: 1 }}>{step.title}</Typography>
                <Typography variant="body2" color="text.secondary">{step.desc}</Typography>
              </Paper>
            </Grid>
          ))}
        </Grid>
      </Container>

      <Box id="raffles-section">
        <Typography variant="h4" sx={{ fontWeight: 700, mb: 3 }}>Rifas Activas</Typography>

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
            {totalPages > 1 && (
              <Pagination count={totalPages} page={page} onChange={(_, p) => setPage(p)} sx={{ mt: 3, display: 'flex', justifyContent: 'center' }} />
            )}
            {raffles.length === 0 && (
              <Box sx={{ textAlign: 'center', py: 8 }}>
                <Typography variant="h6" color="text.secondary" gutterBottom>
                  No hay rifas activas en este momento
                </Typography>
                <Typography color="text.secondary">
                  Vuelve pronto o crea tu propia rifa desde el Dashboard
                </Typography>
              </Box>
            )}
          </>
        )}
      </Box>
    </Box>
  );
}
