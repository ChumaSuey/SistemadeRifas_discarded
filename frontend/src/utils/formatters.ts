export function formatCurrency(amount: number, currency: string): string {
  const symbols: Record<string, string> = { USD: '$', VES: 'Bs.', EUR: '€' };
  const sym = symbols[currency] || currency;
  return `${sym} ${amount.toFixed(2)}`;
}

export function statusLabel(status: string): string {
  const labels: Record<string, string> = {
    draft: 'Borrador',
    pending_verification: 'Pendiente de verificación',
    active: 'Activa',
    finished: 'Finalizada',
    cancelled: 'Cancelada',
    pending: 'Pendiente',
    verified: 'Verificado',
    rejected: 'Rechazado',
    winner: 'Ganador',
    suspended: 'Suspendido',
    banned: 'Baneado',
  };
  return labels[status] || status;
}

export function statusColor(status: string): 'default' | 'success' | 'error' | 'warning' | 'info' {
  const map: Record<string, 'default' | 'success' | 'error' | 'warning' | 'info'> = {
    active: 'success',
    finished: 'info',
    cancelled: 'default',
    draft: 'warning',
    pending: 'warning',
    verified: 'success',
    rejected: 'error',
    winner: 'success',
    suspended: 'warning',
    banned: 'error',
  };
  return map[status] || 'default';
}

export function roleLabel(role: string): string {
  const labels: Record<string, string> = {
    admin: 'Admin',
    raffler: 'Raffler',
    participant: 'Participante',
    verifier: 'Verificador',
  };
  return labels[role] || role;
}
