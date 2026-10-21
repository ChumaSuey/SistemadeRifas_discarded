import { Container, Box } from '@mui/material';
import { Outlet } from 'react-router-dom';
import Navbar from './Navbar';

export default function Layout() {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <Navbar />
      <Container maxWidth="lg" sx={{ mt: { xs: 2, md: 3 }, mb: { xs: 4, md: 6 }, flex: 1 }}>
        <Outlet />
      </Container>
    </Box>
  );
}
