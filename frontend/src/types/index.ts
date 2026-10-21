export interface User {
  id: string;
  name: string;
  email: string;
  email_verified: boolean;
  role: 'admin' | 'raffler' | 'participant' | 'verifier';
  status: 'active' | 'suspended' | 'banned';
  social_links?: {
    twitch?: string;
    instagram?: string;
  };
  created_at: string;
}

export interface Raffle {
  id: string;
  raffler_id: string;
  title: string;
  description: string;
  status: 'draft' | 'pending_verification' | 'active' | 'finished' | 'cancelled';
  pricing_tier: 'free' | 'pro';
  ticket_price: number;
  currency: 'USD' | 'VES' | 'EUR';
  payment_methods: string[];
  tickets_total: number;
  tickets_sold: number;
  winner_id?: string;
  winning_ticket_id?: string;
  sort_commit_hash?: string;
  draw_date?: string;
  created_at: string;
  updated_at: string;
}

export interface Ticket {
  id: string;
  raffle_id: string;
  participant_id: string;
  ticket_number: number;
  status: 'pending' | 'verified' | 'rejected' | 'cancelled' | 'winner';
  payment_method: string;
  payment_proof_url?: string;
  stripe_payment_intent_id?: string;
  rejection_reason?: string;
  verified_by?: string;
  verified_at?: string;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  user: User;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  per_page: number;
  total_pages: number;
}

export interface SortCommitResponse {
  commit_hash: string;
  tickets_count: number;
  timestamp: number;
}

export interface SortRevealResponse {
  commit_hash: string;
  seed: string;
  winning_ticket_id: string;
  winner_user_id: string;
  ticket_number: number;
  draw_date: string;
}

export interface RaffleDetail extends Raffle {
  tickets?: Ticket[];
}
