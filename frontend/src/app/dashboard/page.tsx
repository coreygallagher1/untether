'use client';

import { Typography, Box, Paper, Grid, Button } from '@mui/material';
import { AccountBalance, TrendingUp, Notifications, Assessment } from '@mui/icons-material';
import { useAuth } from '@/app/contexts/AuthContext';

export default function DashboardPage() {
  const { user } = useAuth();

  console.log('Dashboard page - user data:', user);

  const features = [
    {
      icon: <AccountBalance sx={{ fontSize: 40, color: '#1B4D3E' }} />,
      title: 'Financial Accounts',
      description: 'Connect and manage all your bank accounts, credit cards, and investments in one place.',
      action: 'Connect Account'
    },
    {
      icon: <TrendingUp sx={{ fontSize: 40, color: '#1B4D3E' }} />,
      title: 'Budgets & Goals',
      description: 'Set spending limits, create savings goals, and track your progress over time.',
      action: 'Set Budget'
    },
    {
      icon: <Assessment sx={{ fontSize: 40, color: '#1B4D3E' }} />,
      title: 'Financial Insights',
      description: 'Get personalized insights and recommendations based on your spending patterns.',
      action: 'View Insights'
    },
    {
      icon: <Notifications sx={{ fontSize: 40, color: '#1B4D3E' }} />,
      title: 'Alerts & Notifications',
      description: 'Stay informed with customizable alerts for account activity and bill payments.',
      action: 'Set Alerts'
    }
  ];

  const userName = user?.first_name || 'User';

  return (
    <Box sx={{ 
      p: { xs: 2, sm: 3 },
      maxWidth: 1200,
      mx: 'auto'
    }}>
      <Paper 
        elevation={3} 
        sx={{ 
          p: { xs: 3, sm: 4 }, 
          mb: 4,
          bgcolor: '#1B4D3E',
          color: 'white',
          textAlign: 'center'
        }}
      >
        <Typography variant="h4" component="h1" gutterBottom sx={{ 
          fontWeight: 600,
          fontSize: { xs: '1.75rem', sm: '2.125rem' }
        }}>
          Welcome back, {userName}!
        </Typography>
        <Typography variant="body1" sx={{ opacity: 0.9 }}>
          Your financial journey continues here. Let's make your money work smarter.
        </Typography>
      </Paper>

      <Grid container spacing={3}>
        {features.map((feature, index) => (
          <Grid item xs={12} sm={6} key={index}>
            <Paper 
              elevation={2} 
              sx={{ 
                p: 3,
                height: '100%',
                display: 'flex',
                flexDirection: 'column',
                alignItems: 'flex-start',
                transition: 'transform 0.2s ease-in-out, box-shadow 0.2s ease-in-out',
                '&:hover': {
                  transform: 'translateY(-4px)',
                  boxShadow: 4
                }
              }}
            >
              <Box sx={{ mb: 2 }}>
                {feature.icon}
              </Box>
              <Typography variant="h6" gutterBottom sx={{ 
                fontWeight: 600, 
                color: '#1B4D3E',
                width: '100%'
              }}>
                {feature.title}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ 
                mb: 3, 
                flex: 1,
                width: '100%'
              }}>
                {feature.description}
              </Typography>
              <Button 
                variant="outlined" 
                sx={{ 
                  color: '#1B4D3E',
                  borderColor: '#1B4D3E',
                  '&:hover': {
                    borderColor: '#1B4D3E',
                    bgcolor: 'rgba(27, 77, 62, 0.1)'
                  }
                }}
              >
                {feature.action}
              </Button>
            </Paper>
          </Grid>
        ))}
      </Grid>
    </Box>
  );
} 