import { createContext, useContext, useState, useCallback, type ReactNode } from 'react';
import { Snackbar, Alert, type AlertColor } from '@mui/material';

interface SnackbarItem {
  id: number;
  message: string;
  severity: AlertColor;
}

interface SnackbarState {
  showSnackbar: (message: string, severity?: AlertColor) => void;
}

const SnackbarContext = createContext<SnackbarState | null>(null);
let nextId = 0;

export function SnackbarProvider({ children }: { children: ReactNode }) {
  const [items, setItems] = useState<SnackbarItem[]>([]);

  const showSnackbar = useCallback((message: string, severity: AlertColor = 'success') => {
    const id = nextId++;
    setItems((prev) => [...prev, { id, message, severity }]);
    setTimeout(() => {
      setItems((prev) => prev.filter((i) => i.id !== id));
    }, 4000);
  }, []);

  const handleClose = (id: number) => {
    setItems((prev) => prev.filter((i) => i.id !== id));
  };

  return (
    <SnackbarContext.Provider value={{ showSnackbar }}>
      {children}
      {items.map((item, idx) => (
        <Snackbar
          key={item.id}
          open
          autoHideDuration={4000}
          onClose={() => handleClose(item.id)}
          anchorOrigin={{ vertical: 'top', horizontal: 'right' }}
          sx={{ top: idx * 60 + 70 }}
        >
          <Alert
            onClose={() => handleClose(item.id)}
            severity={item.severity}
            variant="filled"
            sx={{ width: '100%' }}
          >
            {item.message}
          </Alert>
        </Snackbar>
      ))}
    </SnackbarContext.Provider>
  );
}

export function useSnackbar() {
  const ctx = useContext(SnackbarContext);
  if (!ctx) throw new Error('useSnackbar must be inside SnackbarProvider');
  return ctx;
}
