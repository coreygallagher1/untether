import React from 'react';
import Box from '@mui/material/Box';
import Typography from '@mui/material/Typography';
import Button from '@mui/material/Button';
import Stack from '@mui/material/Stack';

export default function Home() {
  return (
    <Box
      sx={{
        display: 'flex',
        flexDirection: 'column',
        alignItems: 'center',
        justifyContent: 'center',
        minHeight: 'calc(100vh - 200px)',
        textAlign: 'center',
      }}
    >
      <Stack spacing={4} alignItems="center">
        <Typography variant="h2" component="h1" gutterBottom>
          Welcome to Untether
        </Typography>
        <Typography variant="h5" component="h2" color="text.secondary" gutterBottom>
          Your journey to financial freedom starts here
        </Typography>
        <Button
          variant="contained"
          size="large"
          color="primary"
          sx={{ mt: 4, px: 4, py: 1.5 }}
        >
          Get Started
        </Button>
      </Stack>
    </Box>
  );
}
