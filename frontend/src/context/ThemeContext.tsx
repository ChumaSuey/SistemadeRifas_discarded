import { createContext, useContext, useState, useMemo, useEffect, type ReactNode } from 'react';
import { ThemeProvider, createTheme } from '@mui/material';

interface ThemeState {
  darkMode: boolean;
  toggleDarkMode: () => void;
}

const ThemeContext = createContext<ThemeState | null>(null);

function getStoredMode(): boolean {
  try {
    return localStorage.getItem('darkMode') === 'true';
  } catch {
    return false;
  }
}

export function AppThemeProvider({ children }: { children: ReactNode }) {
  const [darkMode, setDarkMode] = useState(getStoredMode);

  useEffect(() => {
    localStorage.setItem('darkMode', String(darkMode));
  }, [darkMode]);

  const toggleDarkMode = () => setDarkMode((prev) => !prev);

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: darkMode ? 'dark' : 'light',
          primary: { main: darkMode ? '#90caf9' : '#1976d2' },
        },
      }),
    [darkMode],
  );

  return (
    <ThemeContext.Provider value={{ darkMode, toggleDarkMode }}>
      <ThemeProvider theme={theme}>{children}</ThemeProvider>
    </ThemeContext.Provider>
  );
}

export function useThemeMode() {
  const ctx = useContext(ThemeContext);
  if (!ctx) throw new Error('useThemeMode must be inside AppThemeProvider');
  return ctx;
}
