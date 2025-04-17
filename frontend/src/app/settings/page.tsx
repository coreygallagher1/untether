'use client';

import { useAuth } from '@/app/contexts/AuthContext';
import { Box, Paper, Typography, Grid, Skeleton, Avatar, Button } from '@mui/material';
import { AccountCircle, Email, CalendarToday, Edit } from '@mui/icons-material';

export default function SettingsPage() {
  const { user } = useAuth();

  console.log('Settings page - user data:', user);

  if (!user) {
    return (
      <Box sx={{ p: 2, maxWidth: 700, mx: 'auto' }}>
        <Paper elevation={3} sx={{ p: 3 }}>
          <Skeleton variant="text" sx={{ fontSize: '1.5rem', mb: 2 }} />
          <Skeleton variant="rectangular" height={100} />
        </Paper>
      </Box>
    );
  }

  const formatDate = (timestamp: { seconds: number; nanos: number }) => {
    try {
      const date = new Date(timestamp.seconds * 1000 + timestamp.nanos / 1000000);
      return new Intl.DateTimeFormat('en-US', { 
        year: 'numeric', 
        month: 'long', 
        day: 'numeric' 
      }).format(date);
    } catch (error) {
      console.error('Error formatting date:', error);
      return 'Not available';
    }
  };

  return (
    <Box sx={{ p: 2, maxWidth: 700, mx: 'auto' }}>
      <Paper elevation={2} sx={{ borderRadius: 2, overflow: 'hidden' }}>
        {/* Header */}
        <Box sx={{ 
          p: 3,
          bgcolor: '#1B4D3E',
          color: 'white',
          display: 'flex', 
          alignItems: 'center', 
          gap: 2,
        }}>
          <Avatar 
            sx={{ 
              width: 56, 
              height: 56, 
              bgcolor: 'white',
              color: '#1B4D3E',
              fontSize: '1.5rem',
              fontWeight: 600,
              boxShadow: '0 2px 8px rgba(0,0,0,0.15)'
            }}
          >
            {user.first_name.charAt(0)}{user.last_name.charAt(0)}
          </Avatar>
          <Typography variant="h5" component="h1" sx={{ 
            fontWeight: 600
          }}>
            Account Information
          </Typography>
        </Box>

        <Box sx={{ p: 3 }}>
          {/* Account Information Fields */}
          <Grid container spacing={2}>
            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ pl: 1 }}>
                First Name
              </Typography>
              <Box sx={{ 
                p: 2,
                bgcolor: '#F8F9FA',
                borderRadius: 1.5,
                minHeight: '48px',
                display: 'flex',
                alignItems: 'center'
              }}>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                  {user.first_name || 'Not available'}
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12} sm={6}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ pl: 1 }}>
                Last Name
              </Typography>
              <Box sx={{ 
                p: 2,
                bgcolor: '#F8F9FA',
                borderRadius: 1.5,
                minHeight: '48px',
                display: 'flex',
                alignItems: 'center'
              }}>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                  {user.last_name || 'Not available'}
                </Typography>
              </Box>
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ pl: 1 }}>
                Email Address
              </Typography>
              <Box sx={{ 
                p: 2,
                bgcolor: '#F8F9FA',
                borderRadius: 1.5,
                minHeight: '48px',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'space-between'
              }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                  <Email sx={{ color: 'text.secondary', fontSize: 20 }} />
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {user.email || 'Not available'}
                  </Typography>
                </Box>
                <Button 
                  startIcon={<Edit sx={{ fontSize: 18 }} />}
                  variant="outlined"
                  size="small"
                  sx={{ 
                    color: '#1B4D3E',
                    borderColor: '#1B4D3E',
                    py: 0.5,
                    '&:hover': {
                      borderColor: '#1B4D3E',
                      bgcolor: 'rgba(27, 77, 62, 0.1)'
                    }
                  }}
                >
                  Edit
                </Button>
              </Box>
            </Grid>

            <Grid item xs={12}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ pl: 1 }}>
                Member Since
              </Typography>
              <Box sx={{ 
                p: 2,
                bgcolor: '#F8F9FA',
                borderRadius: 1.5,
                minHeight: '48px',
                display: 'flex',
                alignItems: 'center'
              }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 1.5 }}>
                  <CalendarToday sx={{ color: 'text.secondary', fontSize: 20 }} />
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {formatDate(user.created_at)}
                  </Typography>
                </Box>
              </Box>
            </Grid>
          </Grid>

          {/* Account Actions Section */}
          <Box sx={{ mt: 4 }}>
            <Typography variant="h6" gutterBottom sx={{ 
              color: '#1B4D3E', 
              fontWeight: 600,
              mb: 2
            }}>
              Account Actions
            </Typography>
            <Grid container spacing={2}>
              <Grid item xs={12} sm={6}>
                <Button
                  fullWidth
                  variant="outlined"
                  color="error"
                  size="large"
                  sx={{ 
                    borderColor: 'error.main',
                    borderRadius: 1.5,
                    '&:hover': {
                      borderColor: 'error.dark',
                      bgcolor: 'error.light',
                      opacity: 0.1
                    }
                  }}
                >
                  Delete Account
                </Button>
              </Grid>
              <Grid item xs={12} sm={6}>
                <Button
                  fullWidth
                  variant="contained"
                  size="large"
                  sx={{ 
                    bgcolor: '#1B4D3E',
                    borderRadius: 1.5,
                    '&:hover': {
                      bgcolor: '#143C30'
                    }
                  }}
                >
                  Change Password
                </Button>
              </Grid>
            </Grid>
          </Box>
        </Box>
      </Paper>
    </Box>
  );
} 