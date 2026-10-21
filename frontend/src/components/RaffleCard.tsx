import { useState, useEffect } from 'react';
import { Card, CardContent, Typography, Chip, CardActionArea, Box } from '@mui/material';
import AccessTimeIcon from '@mui/icons-material/AccessTime';
import { useNavigate } from 'react-router-dom';
import type { Raffle } from '../types';
import { formatCurrency, statusLabel, statusColor } from '../utils/formatters';

function Countdown({ drawDate }: { drawDate: string }) {
  const [now, setNow] = useState(Date.now());

  useEffect(() => {
    const id = setInterval(() => setNow(Date.now()), 1000);
    return () => clearInterval(id);
  }, []);

  const target = new Date(drawDate).getTime();
  const diff = target - now;

  if (diff <= 0) {
    return (
      <Typography variant="caption" color="error" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
        <AccessTimeIcon fontSize="inherit" /> Finalizado
      </Typography>
    );
  }

  const days = Math.floor(diff / 86400000);
  const hours = Math.floor((diff % 86400000) / 3600000);
  const mins = Math.floor((diff % 3600000) / 60000);

  const label = days > 0
    ? `${days}d ${hours}h`
    : hours > 0
      ? `${hours}h ${mins}m`
      : `${mins}m`;

  return (
    <Typography variant="caption" color="warning.main" sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
      <AccessTimeIcon fontSize="inherit" /> {label}
    </Typography>
  );
}

export default function RaffleCard({ raffle }: { raffle: Raffle }) {
  const navigate = useNavigate();

  return (
    <Card
      sx={{
        transition: 'transform 0.2s, box-shadow 0.2s',
        '&:hover': { transform: 'translateY(-3px)', boxShadow: 6 },
      }}
    >
      <CardActionArea onClick={() => navigate(`/raffles/${raffle.id}`)}>
        <CardContent>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'flex-start', mb: 1 }}>
            <Typography variant="h6" sx={{ flex: 1, fontSize: { xs: '1rem', sm: '1.1rem', md: '1.25rem' } }}>{raffle.title}</Typography>
            <Chip label={statusLabel(raffle.status)} color={statusColor(raffle.status)} size="small" />
          </Box>
          <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
            {raffle.description.length > 100 ? raffle.description.slice(0, 100) + '...' : raffle.description}
          </Typography>
          <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <Box sx={{ display: 'flex', gap: 2 }}>
              <Typography variant="body2">{formatCurrency(raffle.ticket_price, raffle.currency)}</Typography>
              <Typography variant="body2" color="text.secondary">
                {raffle.tickets_sold}/{raffle.tickets_total} boletos
              </Typography>
            </Box>
            {raffle.draw_date && raffle.status === 'active' && (
              <Countdown drawDate={raffle.draw_date} />
            )}
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
