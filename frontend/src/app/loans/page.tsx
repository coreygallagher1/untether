'use client';

import { Box, Typography, Paper, Grid, Card, CardContent, CardActions, Button, Chip, LinearProgress, Alert } from '@mui/material';
import { CreditCard, Add, TrendingUp, Payment, Schedule, CheckCircle } from '@mui/icons-material';

export default function LoansPage() {
  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: 'auto' }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ 
          fontWeight: 600, 
          color: '#1B4D3E',
          mb: 1
        }}>
          Loans & Credit
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Manage your loans, track payments, and explore refinancing options
        </Typography>
      </Box>

      {/* Quick Actions */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
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
          Add New Loan
        </Button>
        <Button
          variant="outlined"
          startIcon={<TrendingUp />}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Refinance Options
        </Button>
        <Button
          variant="outlined"
          startIcon={<Payment />}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Payment Calculator
        </Button>
      </Box>

      {/* Alert */}
      <Alert severity="info" sx={{ mb: 3 }}>
        <Typography variant="body2">
          <strong>Tip:</strong> Making extra payments or paying bi-weekly can significantly reduce your total interest paid and help you pay off loans faster.
        </Typography>
      </Alert>

      {/* Loans Grid */}
      <Grid container spacing={3}>
        {/* Sample Loan 1 - Student Loan */}
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
                <CreditCard sx={{ fontSize: 32, color: '#1B4D3E', mr: 2 }} />
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Student Loan
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Federal Direct Loan
                  </Typography>
                </Box>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="h4" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $18,247.32
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Remaining Balance
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    Progress
                  </Typography>
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    $6,752 / $25,000
                  </Typography>
                </Box>
                <LinearProgress 
                  variant="determinate" 
                  value={27} 
                  sx={{ 
                    height: 8, 
                    borderRadius: 4,
                    bgcolor: '#e0e0e0',
                    '& .MuiLinearProgress-bar': {
                      bgcolor: '#1B4D3E'
                    }
                  }} 
                />
              </Box>

              <Grid container spacing={2} sx={{ mb: 2 }}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Interest Rate
                  </Typography>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    4.25%
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Monthly Payment
                  </Typography>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    $247.50
                  </Typography>
                </Grid>
              </Grid>

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="Federal" 
                  color="primary" 
                  size="small" 
                />
              </Box>

              <Typography variant="body2" color="text.secondary">
                Next payment due: Dec 15, 2024
              </Typography>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
                Make Payment
              </Button>
            </CardActions>
          </Card>
        </Grid>

        {/* Sample Loan 2 - Auto Loan */}
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
                <CreditCard sx={{ fontSize: 32, color: '#1B4D3E', mr: 2 }} />
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Auto Loan
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    2022 Honda Civic
                  </Typography>
                </Box>
              </Box>
              
              <Box sx={{ mb: 2 }}>
                <Typography variant="h4" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $12,847.67
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Remaining Balance
                </Typography>
              </Box>

              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    Progress
                  </Typography>
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    $7,152 / $20,000
                  </Typography>
                </Box>
                <LinearProgress 
                  variant="determinate" 
                  value={36} 
                  sx={{ 
                    height: 8, 
                    borderRadius: 4,
                    bgcolor: '#e0e0e0',
                    '& .MuiLinearProgress-bar': {
                      bgcolor: '#1B4D3E'
                    }
                  }} 
                />
              </Box>

              <Grid container spacing={2} sx={{ mb: 2 }}>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Interest Rate
                  </Typography>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    3.89%
                  </Typography>
                </Grid>
                <Grid item xs={6}>
                  <Typography variant="body2" color="text.secondary">
                    Monthly Payment
                  </Typography>
                  <Typography variant="h6" sx={{ fontWeight: 600 }}>
                    $387.25
                  </Typography>
                </Grid>
              </Grid>

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="Private" 
                  color="secondary" 
                  size="small" 
                />
              </Box>

              <Typography variant="body2" color="text.secondary">
                Next payment due: Dec 20, 2024
              </Typography>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
                Make Payment
              </Button>
            </CardActions>
          </Card>
        </Grid>
      </Grid>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mt: 2 }}>
        {/* Total Debt Summary */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, bgcolor: '#f8f9fa' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <TrendingUp sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Total Debt Summary
              </Typography>
            </Box>
            
            <Grid container spacing={2}>
              <Grid item xs={6}>
                <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $31,094.99
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Total Debt
                </Typography>
              </Grid>
              <Grid item xs={6}>
                <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                  $634.75
                </Typography>
                <Typography variant="body2" color="text.secondary">
                  Monthly Payments
                </Typography>
              </Grid>
            </Grid>
          </Paper>
        </Grid>

        {/* Payment Strategy */}
        <Grid item xs={12} md={6}>
          <Paper sx={{ p: 3, bgcolor: '#f8f9fa' }}>
            <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
              <Schedule sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
              <Typography variant="h6" sx={{ fontWeight: 600 }}>
                Payment Strategy
              </Typography>
            </Box>
            
            <Box sx={{ mb: 2 }}>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                Recommended: Snowball Method
              </Typography>
              <Typography variant="body2">
                Focus on paying off the auto loan first, then apply extra payments to the student loan.
              </Typography>
            </Box>

            <Box sx={{ display: 'flex', alignItems: 'center', gap: 1 }}>
              <CheckCircle sx={{ fontSize: 16, color: 'success.main' }} />
              <Typography variant="body2" color="text.secondary">
                Could save $2,847 in interest with current strategy
              </Typography>
            </Box>
          </Paper>
        </Grid>
      </Grid>
    </Box>
  );
}
