'use client';

import { Box, Typography, Paper, Grid, Card, CardContent, CardActions, Button, Chip } from '@mui/material';
import { AccountBalance, Add, TrendingUp } from '@mui/icons-material';

export default function BankAccountsPage() {
  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: 'auto' }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ 
          fontWeight: 600, 
          color: '#1B4D3E',
          mb: 1
        }}>
          Bank Accounts
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Manage your connected bank accounts and view your financial overview
        </Typography>
      </Box>

      {/* Add Account Button */}
      <Box sx={{ mb: 3 }}>
        <Button
          variant="contained"
          startIcon={<Add />}
          sx={{
            bgcolor: '#1B4D3E',
            '&:hover': {
              bgcolor: '#143C30'
            }
          }}
        >
          Connect Bank Account
        </Button>
      </Box>

      {/* Accounts Grid */}
      <Grid container spacing={3}>
        {/* Sample Account 1 */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            height: '100%',
            border: '1px solid #e0e0e0',
            '&:hover': {
              boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
              transform: 'translateY(-2px)',
              transition: 'all 0.2s ease-in-out'
            }
          }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AccountBalance sx={{ fontSize: 32, color: '#1B4D3E', mr: 2 }} />
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Test Checking Account
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    ****1234
                  </Typography>
                </Box>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="h4" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $2,847.32
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Available Balance
                </Typography>
              </Box>

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="Roundup Enabled" 
                  color="primary" 
                  size="small" 
                />
              </Box>

              <Typography variant="body2" color="text.secondary">
                Last updated: 2 hours ago
              </Typography>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" color="secondary">
                Settings
              </Button>
            </CardActions>
          </Card>
        </Grid>

        {/* Sample Account 2 */}
        <Grid item xs={12} md={6}>
          <Card sx={{ 
            height: '100%',
            border: '1px solid #e0e0e0',
            '&:hover': {
              boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
              transform: 'translateY(-2px)',
              transition: 'all 0.2s ease-in-out'
            }
          }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AccountBalance sx={{ fontSize: 32, color: '#1B4D3E', mr: 2 }} />
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Test Savings Account
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    ****5678
                  </Typography>
                </Box>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="h4" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $8,421.67
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Available Balance
                </Typography>
              </Box>

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="Roundup Enabled" 
                  color="primary" 
                  size="small" 
                />
              </Box>

              <Typography variant="body2" color="text.secondary">
                Last updated: 1 hour ago
              </Typography>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" color="secondary">
                Settings
              </Button>
            </CardActions>
          </Card>
        </Grid>
      </Grid>

      {/* Summary Card */}
      <Paper sx={{ mt: 4, p: 3, bgcolor: '#f8f9fa' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <TrendingUp sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            Account Summary
          </Typography>
        </Box>
        
        <Grid container spacing={3}>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              $11,268.99
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Total Balance
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              $47.23
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Total Roundups This Month
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              2
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Connected Accounts
            </Typography>
          </Grid>
        </Grid>
      </Paper>
    </Box>
  );
}
