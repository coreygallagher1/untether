'use client';

import { createTheme } from '@mui/material/styles';

export const theme = createTheme({
  palette: {
    primary: {
      main: '#1B4D3E',
      light: '#2C7A62',
      dark: '#0B2318',
      contrastText: '#ffffff',
    },
    secondary: {
      main: '#38A169',
      light: '#48BB78',
      dark: '#2F855A',
      contrastText: '#ffffff',
    },
    background: {
      default: '#f8f9fa',
      paper: '#ffffff',
    },
    text: {
      primary: '#2D3748',
      secondary: '#4A5568',
    },
  },
  components: {
    MuiButton: {
      styleOverrides: {
        root: {
          textTransform: 'none',
          borderRadius: 8,
        },
      },
    },
    MuiAppBar: {
      styleOverrides: {
        root: {
          backgroundColor: '#1B4D3E',
        },
      },
    },
  },
}); 