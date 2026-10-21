import { useState } from 'react';
import { AppBar, Toolbar, Typography, Button, Box, Chip, IconButton, Drawer, List, ListItem, ListItemButton, Divider } from '@mui/material';
import MenuIcon from '@mui/icons-material/Menu';
import DarkModeIcon from '@mui/icons-material/DarkMode';
import LightModeIcon from '@mui/icons-material/LightMode';
import { Link as RouterLink, useNavigate } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';
import { useThemeMode } from '../context/ThemeContext';
import { roleLabel } from '../utils/formatters';

export default function Navbar() {
  const { user, isAuthenticated, isRaffler, isVerifier, isAdmin, logout } = useAuth();
  const { darkMode, toggleDarkMode } = useThemeMode();
  const navigate = useNavigate();
  const [drawerOpen, setDrawerOpen] = useState(false);

  const handleLogout = async () => {
    setDrawerOpen(false);
    await logout();
    navigate('/');
  };

  const navLinks = isAuthenticated
    ? [
        ...(isRaffler ? [{ label: 'Dashboard', to: '/dashboard' }] : []),
        ...(isVerifier ? [{ label: 'Verificar', to: '/verifier' }] : []),
        { label: 'Mis Tickets', to: '/tickets/my' },
        ...(isAdmin ? [{ label: 'Admin', to: '/admin' }] : []),
      ]
    : [{ label: 'Entrar', to: '/login' }];

  return (
    <>
      <AppBar position="sticky">
        <Toolbar>
          <Typography
            variant="h6"
            component={RouterLink}
            to="/"
            sx={{ textDecoration: 'none', color: 'inherit', flexGrow: 1, fontSize: { xs: '1.1rem', md: '1.25rem' } }}
          >
            Sistema de Rifas
          </Typography>

          <Box sx={{ display: { xs: 'none', md: 'flex' }, gap: 1, alignItems: 'center' }}>
            {isAuthenticated ? (
              <>
                {navLinks.map((link) => (
                  <Button key={link.to} color="inherit" component={RouterLink} to={link.to} sx={{ fontSize: '0.85rem' }}>
                    {link.label}
                  </Button>
                ))}
                <Chip label={roleLabel(user!.role)} size="small" color="secondary" />
                <IconButton color="inherit" onClick={toggleDarkMode} size="small">
                  {darkMode ? <LightModeIcon /> : <DarkModeIcon />}
                </IconButton>
                <Button color="inherit" onClick={handleLogout} sx={{ fontSize: '0.85rem' }}>Salir</Button>
              </>
            ) : (
              <Button color="inherit" component={RouterLink} to="/login">Entrar</Button>
            )}
          </Box>

          <IconButton
            color="inherit"
            edge="end"
            onClick={() => setDrawerOpen(true)}
            sx={{ display: { xs: 'flex', md: 'none' } }}
          >
            <MenuIcon />
          </IconButton>
        </Toolbar>
      </AppBar>

      <Drawer anchor="right" open={drawerOpen} onClose={() => setDrawerOpen(false)}>
        <Box sx={{ width: 260, pt: 2 }}>
          {isAuthenticated && user && (
            <Box sx={{ px: 2, pb: 2 }}>
              <Typography variant="subtitle1" sx={{ fontWeight: 700 }}>{user.name}</Typography>
              <Chip label={roleLabel(user.role)} size="small" color="secondary" sx={{ mt: 0.5 }} />
            </Box>
          )}
          <Divider />
          <List>
            <ListItem disablePadding>
              <ListItemButton component={RouterLink} to="/" onClick={() => setDrawerOpen(false)}>
                Inicio
              </ListItemButton>
            </ListItem>
            {navLinks.map((link) => (
              <ListItem key={link.to} disablePadding>
                <ListItemButton component={RouterLink} to={link.to} onClick={() => setDrawerOpen(false)}>
                  {link.label}
                </ListItemButton>
              </ListItem>
            ))}
          </List>
          {isAuthenticated && (
            <>
              <Divider />
              <List>
                <ListItem disablePadding>
                  <ListItemButton onClick={() => { toggleDarkMode(); setDrawerOpen(false); }}>
                    {darkMode ? 'Modo Claro' : 'Modo Oscuro'}
                  </ListItemButton>
                </ListItem>
                <ListItem disablePadding>
                  <ListItemButton onClick={handleLogout}>Salir</ListItemButton>
                </ListItem>
              </List>
            </>
          )}
        </Box>
      </Drawer>
    </>
  );
}
