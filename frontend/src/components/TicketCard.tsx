import { Card, CardContent, Typography, Chip, Box, CardActionArea } from '@mui/material';
import { useNavigate } from 'react-router-dom';
import type { Ticket } from '../types';
import { formatCurrency, statusLabel, statusColor } from '../utils/formatters';

export default function TicketCard({ ticket, price, currency }: { ticket: Ticket; price?: number; currency?: string }) {
  const navigate = useNavigate();

  return (
    <Card variant="outlined" sx={{ mb: 1 }}>
      <CardActionArea onClick={() => navigate(`/raffles/${ticket.raffle_id}`)}>
        <CardContent sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', py: 1.5, '&:last-child': { pb: 1.5 } }}>
          <Box>
            <Typography variant="subtitle1">Ticket #{ticket.ticket_number}</Typography>
            <Typography variant="body2" color="text.secondary">
              {ticket.payment_method === 'manual_transfer' ? 'Transferencia' : 'Stripe'}
            </Typography>
          </Box>
          <Box sx={{ display: 'flex', gap: 1, alignItems: 'center' }}>
            {price && currency && (
              <Typography variant="body2">{formatCurrency(price, currency)}</Typography>
            )}
            <Chip label={statusLabel(ticket.status)} color={statusColor(ticket.status)} size="small" />
            {ticket.rejection_reason && (
              <Typography variant="caption" color="error">{ticket.rejection_reason}</Typography>
            )}
          </Box>
        </CardContent>
      </CardActionArea>
    </Card>
  );
}
