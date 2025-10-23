'use client';

import { Box, Typography, Paper, Grid, Card, CardContent, CardActions, Button, Chip, Avatar, LinearProgress } from '@mui/material';
import { Favorite, Add, TrendingUp, People, AttachMoney } from '@mui/icons-material';

export default function CausesPage() {
  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: 'auto' }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ 
          fontWeight: 600, 
          color: '#1B4D3E',
          mb: 1
        }}>
          Causes
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Support causes you care about with your roundup donations
        </Typography>
      </Box>

      {/* Add Cause Button */}
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
          Add New Cause
        </Button>
      </Box>

      {/* Causes Grid */}
      <Grid container spacing={3}>
        {/* Sample Cause 1 */}
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
                <Avatar sx={{ 
                  bgcolor: '#ff6b6b', 
                  width: 48, 
                  height: 48,
                  mr: 2
                }}>
                  <Favorite />
                </Avatar>
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Local Food Bank
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Fighting hunger in our community
                  </Typography>
                </Box>
              </Box>
              
              <Typography variant="body2" sx={{ mb: 2 }}>
                Help provide meals for families in need. Your roundup donations go directly to purchasing fresh food and essential supplies.
              </Typography>

              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    Progress
                  </Typography>
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    $2,847 / $5,000
                  </Typography>
                </Box>
                <LinearProgress 
                  variant="determinate" 
                  value={57} 
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

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="Local" 
                  color="primary" 
                  size="small" 
                />
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <People sx={{ fontSize: 16, color: 'text.secondary' }} />
                  <Typography variant="body2" color="text.secondary">
                    127 supporters
                  </Typography>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <AttachMoney sx={{ fontSize: 16, color: 'text.secondary' }} />
                  <Typography variant="body2" color="text.secondary">
                    $23.50 donated
                  </Typography>
                </Box>
              </Box>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
                Donate Now
              </Button>
            </CardActions>
          </Card>
        </Grid>

        {/* Sample Cause 2 */}
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
                <Avatar sx={{ 
                  bgcolor: '#4ecdc4', 
                  width: 48, 
                  height: 48,
                  mr: 2
                }}>
                  <Favorite />
                </Avatar>
                <Box>
                  <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                    Education Fund
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Supporting students' educational needs
                  </Typography>
                </Box>
              </Box>
              
              <Typography variant="body2" sx={{ mb: 2 }}>
                Help provide scholarships, school supplies, and educational resources for students in underserved communities.
              </Typography>

              <Box sx={{ mb: 2 }}>
                <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                  <Typography variant="body2" color="text.secondary">
                    Progress
                  </Typography>
                  <Typography variant="body2" sx={{ fontWeight: 600 }}>
                    $8,421 / $10,000
                  </Typography>
                </Box>
                <LinearProgress 
                  variant="determinate" 
                  value={84} 
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

              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <Chip 
                  label="Active" 
                  color="success" 
                  size="small" 
                />
                <Chip 
                  label="National" 
                  color="secondary" 
                  size="small" 
                />
              </Box>

              <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <People sx={{ fontSize: 16, color: 'text.secondary' }} />
                  <Typography variant="body2" color="text.secondary">
                    89 supporters
                  </Typography>
                </Box>
                <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                  <AttachMoney sx={{ fontSize: 16, color: 'text.secondary' }} />
                  <Typography variant="body2" color="text.secondary">
                    $47.23 donated
                  </Typography>
                </Box>
              </Box>
            </CardContent>
            
            <CardActions>
              <Button size="small" color="primary">
                View Details
              </Button>
              <Button size="small" variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
                Donate Now
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
            Your Impact Summary
          </Typography>
        </Box>
        
        <Grid container spacing={3}>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              $70.73
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Total Donated
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              2
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Active Causes
            </Typography>
          </Grid>
          <Grid item xs={12} sm={4}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              216
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Lives Impacted
            </Typography>
          </Grid>
        </Grid>
      </Paper>
    </Box>
  );
}
