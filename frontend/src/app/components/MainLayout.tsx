'use client';

import React from 'react';
import AppBar from '@mui/material/AppBar';
import Toolbar from '@mui/material/Toolbar';
import Typography from '@mui/material/Typography';
import IconButton from '@mui/material/IconButton';
import MenuIcon from '@mui/icons-material/Menu';
import Box from '@mui/material/Box';
import Container from '@mui/material/Container';
import Image from 'next/image';

export function MainLayout({ children }: { children: React.ReactNode }) {
  return (
    <Box sx={{ display: 'flex', flexDirection: 'column', minHeight: '100vh' }}>
      <AppBar position="static" sx={{ minHeight: 80 }}>
        <Toolbar sx={{ minHeight: 80, py: 2 }}>
          <IconButton
            size="large"
            edge="start"
            color="inherit"
            aria-label="menu"
            sx={{ mr: 2 }}
          >
            <MenuIcon />
          </IconButton>
          <Box sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: 2,
            '& img': {
              filter: 'drop-shadow(0 0 5px rgba(255, 255, 255, 0.6))',
              transition: 'filter 0.3s ease',
              '&:hover': {
                filter: 'drop-shadow(0 0 25px rgba(255, 255, 255, 0.8))',
              }
            }
          }}>
            <Image
              src="/assets/UntetherLogo.png"
              alt="Untether Logo"
              width={56}
              height={56}
              style={{ objectFit: 'contain' }}
            />
            <Typography 
              variant="h6" 
              component="div" 
              sx={{ 
                flexGrow: 1, 
                fontSize: '1.75rem', 
                fontWeight: 600,
                fontFamily: 'var(--font-space-grotesk)',
                letterSpacing: '0.5px'
              }}
            >
              Untether
            </Typography>
          </Box>
        </Toolbar>
      </AppBar>
      <Container component="main" sx={{ mt: 4, mb: 4, flex: 1 }}>
        {children}
      </Container>
    </Box>
  );
} 