import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import api from '../api/client';
import { getStoredUser, setStoredUser, clearTokens, setTokens } from '../utils/storage';
import type { User, AuthResponse } from '../types';

interface AuthState {
  user: User | null;
  loading: boolean;
  login: (email: string, password: string) => Promise<void>;
  register: (name: string, email: string, password: string) => Promise<void>;
  logout: () => Promise<void>;
  isAuthenticated: boolean;
  isRaffler: boolean;
  isVerifier: boolean;
  isAdmin: boolean;
}

const AuthContext = createContext<AuthState | null>(null);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(getStoredUser());
  const [loading, setLoading] = useState(false);

  const authenticate = useCallback((data: AuthResponse) => {
    setTokens(data.access_token, data.refresh_token);
    setStoredUser(data.user);
    setUser(data.user);
  }, []);

  const login = useCallback(async (email: string, password: string) => {
    setLoading(true);
    try {
      const { data } = await api.post<AuthResponse>('/auth/login', { email, password });
      authenticate(data);
    } finally {
      setLoading(false);
    }
  }, [authenticate]);

  const register = useCallback(async (name: string, email: string, password: string) => {
    setLoading(true);
    try {
      const { data } = await api.post<AuthResponse>('/auth/register', { name, email, password });
      authenticate(data);
    } finally {
      setLoading(false);
    }
  }, [authenticate]);

  const logout = useCallback(async () => {
    try {
      await api.post('/auth/logout', { refresh_token: localStorage.getItem('refresh_token') });
    } catch {
      // ignore
    }
    clearTokens();
    setUser(null);
  }, []);

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        register,
        logout,
        isAuthenticated: !!user,
        isRaffler: user?.role === 'raffler' || user?.role === 'admin',
        isVerifier: user?.role === 'verifier' || user?.role === 'admin',
        isAdmin: user?.role === 'admin',
      }}
    >
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const ctx = useContext(AuthContext);
  if (!ctx) throw new Error('useAuth must be inside AuthProvider');
  return ctx;
}
