'use client';

import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  Card, 
  CardContent, 
  CardActions, 
  Button, 
  Chip, 
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  Alert,
  IconButton,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Switch,
  FormControlLabel,
  LinearProgress,
  Avatar
} from '@mui/material';
import { 
  CreditCard, 
  AccountBalance, 
  Add, 
  CalendarToday, 
  Payment, 
  Warning, 
  CheckCircle,
  Schedule,
  AttachMoney,
  TrendingUp,
  TrendingDown,
  Star,
  AutoAwesome,
  Security,
  Speed,
  Savings
} from '@mui/icons-material';
import { useState } from 'react';

export default function BillsPage() {
  const [addBillOpen, setAddBillOpen] = useState(false);
  const [untetherEnabled, setUntetherEnabled] = useState(true);

  // Sample data for Untether Credit Card Rewards system
  const untetherStats = {
    totalRewardsEarned: 1247.50,
    monthlyRewards: 89.32,
    billsPaidThisMonth: 12,
    averageRewardRate: 2.1,
    totalSavings: 1247.50
  };

  const bills = [
    {
      id: 1,
      name: 'Rent Payment',
      amount: 1850.00,
      dueDate: '2024-12-01',
      category: 'Housing',
      status: 'paid',
      untetherEnabled: true,
      rewardsEarned: 37.00,
      paymentMethod: 'Untether Credit Card',
      lastPaid: '2024-11-01'
    },
    {
      id: 2,
      name: 'Electric Bill',
      amount: 125.50,
      dueDate: '2024-12-15',
      category: 'Utilities',
      status: 'upcoming',
      untetherEnabled: true,
      rewardsEarned: 2.51,
      paymentMethod: 'Untether Credit Card',
      lastPaid: '2024-11-15'
    },
    {
      id: 3,
      name: 'Internet Service',
      amount: 79.99,
      dueDate: '2024-12-20',
      category: 'Utilities',
      status: 'upcoming',
      untetherEnabled: true,
      rewardsEarned: 1.60,
      paymentMethod: 'Untether Credit Card',
      lastPaid: '2024-11-20'
    },
    {
      id: 4,
      name: 'Car Insurance',
      amount: 245.00,
      dueDate: '2024-12-05',
      category: 'Insurance',
      status: 'upcoming',
      untetherEnabled: true,
      rewardsEarned: 4.90,
      paymentMethod: 'Untether Credit Card',
      lastPaid: '2024-11-05'
    },
    {
      id: 5,
      name: 'Student Loan',
      amount: 247.50,
      dueDate: '2024-12-15',
      category: 'Loans',
      status: 'upcoming',
      untetherEnabled: false,
      rewardsEarned: 0,
      paymentMethod: 'Direct Debit',
      lastPaid: '2024-11-15'
    },
    {
      id: 6,
      name: 'Phone Bill',
      amount: 95.00,
      dueDate: '2024-12-10',
      category: 'Utilities',
      status: 'upcoming',
      untetherEnabled: true,
      rewardsEarned: 1.90,
      paymentMethod: 'Untether Credit Card',
      lastPaid: '2024-11-10'
    }
  ];

  const upcomingBills = bills.filter(bill => bill.status === 'upcoming');
  const paidBills = bills.filter(bill => bill.status === 'paid');
  const untetherBills = bills.filter(bill => bill.untetherEnabled);
  const totalUpcoming = upcomingBills.reduce((sum, bill) => sum + bill.amount, 0);
  const totalRewardsThisMonth = bills.reduce((sum, bill) => sum + bill.rewardsEarned, 0);

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric'
    });
  };

  const getDaysUntilDue = (dueDate: string) => {
    const today = new Date();
    const due = new Date(dueDate);
    const diffTime = due.getTime() - today.getTime();
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
    return diffDays;
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'overdue': return 'error';
      case 'upcoming': return 'warning';
      case 'paid': return 'success';
      default: return 'default';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'overdue': return <Warning sx={{ color: 'error.main' }} />;
      case 'upcoming': return <Schedule sx={{ color: 'warning.main' }} />;
      case 'paid': return <CheckCircle sx={{ color: 'success.main' }} />;
      default: return <Schedule />;
    }
  };

  const handlePayBill = (billId: number) => {
    console.log(`Paying bill ${billId} with Untether Credit Card`);
  };

  const handleAddBill = () => {
    setAddBillOpen(true);
  };

  const handleCloseAddBill = () => {
    setAddBillOpen(false);
  };

  const toggleUntether = (billId: number) => {
    console.log(`Toggling Untether for bill ${billId}`);
  };

  return (
    <Box sx={{ p: 3, maxWidth: 1200, mx: 'auto' }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ 
          fontWeight: 600, 
          color: '#1B4D3E',
          mb: 1
        }}>
          Smart Bill Pay
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Earn credit card rewards on all your bills with automated payments
        </Typography>
      </Box>

      {/* Untether Credit Card Rewards Banner */}
      <Paper sx={{ 
        p: 3, 
        mb: 4, 
        bgcolor: '#1B4D3E', 
        color: 'white',
        background: 'linear-gradient(135deg, #1B4D3E 0%, #2E7D32 100%)'
      }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h5" sx={{ fontWeight: 600, mb: 1 }}>
              🎉 Untether Credit Card Rewards
            </Typography>
            <Typography variant="body1" sx={{ opacity: 0.9, mb: 2 }}>
              Earn {untetherStats.averageRewardRate}% cashback on all your bills automatically
            </Typography>
            <Box sx={{ display: 'flex', gap: 3 }}>
              <Box>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {formatCurrency(untetherStats.totalRewardsEarned)}
                </Typography>
                <Typography variant="body2" sx={{ opacity: 0.8 }}>
                  Total Rewards Earned
                </Typography>
              </Box>
              <Box>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {formatCurrency(untetherStats.monthlyRewards)}
                </Typography>
                <Typography variant="body2" sx={{ opacity: 0.8 }}>
                  This Month
                </Typography>
              </Box>
              <Box>
                <Typography variant="h6" sx={{ fontWeight: 600 }}>
                  {untetherStats.billsPaidThisMonth}
                </Typography>
                <Typography variant="body2" sx={{ opacity: 0.8 }}>
                  Bills Paid
                </Typography>
              </Box>
            </Box>
          </Box>
          <Box sx={{ textAlign: 'center' }}>
            <Avatar sx={{ 
              bgcolor: 'white', 
              color: '#1B4D3E', 
              width: 80, 
              height: 80,
              mb: 1
            }}>
              <CreditCard sx={{ fontSize: 40 }} />
            </Avatar>
            <Typography variant="body2" sx={{ opacity: 0.9 }}>
              Untether Credit Card
            </Typography>
          </Box>
        </Box>
      </Paper>

      {/* Quick Actions */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
        <Button
          variant="contained"
          startIcon={<Add />}
          onClick={handleAddBill}
          sx={{
            bgcolor: '#1B4D3E',
            '&:hover': {
              bgcolor: '#143C30'
            }
          }}
        >
          Add New Bill
        </Button>
        <Button
          variant="outlined"
          startIcon={<AutoAwesome />}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Enable All Untether
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
          Pay All Bills
        </Button>
      </Box>

      {/* Benefits Alert */}
      <Alert severity="success" sx={{ mb: 3 }}>
        <Typography variant="body2">
          <strong>Smart Savings:</strong> You've earned {formatCurrency(totalRewardsThisMonth)} in rewards this month by using Untether Credit Card for bill payments!
        </Typography>
      </Alert>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Schedule sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Upcoming Bills
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {upcomingBills.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {formatCurrency(totalUpcoming)} total
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Star sx={{ fontSize: 24, color: '#f57c00', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Rewards This Month
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#f57c00' }}>
                {formatCurrency(totalRewardsThisMonth)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {untetherBills.length} bills earning rewards
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AutoAwesome sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Untether Enabled
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {untetherBills.length}/{bills.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Bills earning rewards
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <CheckCircle sx={{ fontSize: 24, color: 'success.main', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Paid This Month
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: 'success.main' }}>
                {paidBills.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Bills paid automatically
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Bills List */}
      <Typography variant="h5" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
        Your Bills
      </Typography>

      <Grid container spacing={3}>
        {bills.map((bill) => (
          <Grid item xs={12} md={6} key={bill.id}>
            <Card sx={{ 
              height: '100%',
              border: bill.untetherEnabled ? '2px solid #1B4D3E' : '1px solid #e0e0e0',
              '&:hover': {
                boxShadow: '0 4px 12px rgba(0,0,0,0.1)',
                transform: 'translateY(-2px)',
                transition: 'all 0.2s ease-in-out'
              }
            }}>
              <CardContent>
                <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                  <Box sx={{ mr: 2 }}>
                    {bill.untetherEnabled ? (
                      <Avatar sx={{ bgcolor: '#1B4D3E', width: 40, height: 40 }}>
                        <CreditCard sx={{ fontSize: 20 }} />
                      </Avatar>
                    ) : (
                      <Avatar sx={{ bgcolor: 'grey.300', width: 40, height: 40 }}>
                        <AccountBalance sx={{ fontSize: 20 }} />
                      </Avatar>
                    )}
                  </Box>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                      {bill.name}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {bill.category} • {bill.paymentMethod}
                    </Typography>
                  </Box>
                  <FormControlLabel
                    control={
                      <Switch
                        checked={bill.untetherEnabled}
                        onChange={() => toggleUntether(bill.id)}
                        color="success"
                      />
                    }
                    label=""
                  />
                </Box>
                
                <Box sx={{ mb: 2 }}>
                  <Typography variant="h4" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                    {formatCurrency(bill.amount)}
                  </Typography>
                  <Typography variant="body2" color="text.secondary">
                    Due {formatDate(bill.dueDate)}
                  </Typography>
                </Box>

                {bill.untetherEnabled && (
                  <Box sx={{ mb: 2, p: 2, bgcolor: '#f8f9fa', borderRadius: 1 }}>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 1 }}>
                      <Star sx={{ fontSize: 16, color: '#f57c00', mr: 1 }} />
                      <Typography variant="body2" sx={{ fontWeight: 600 }}>
                        Rewards Earned: {formatCurrency(bill.rewardsEarned)}
                      </Typography>
                    </Box>
                    <Typography variant="body2" color="text.secondary">
                      Earns {untetherStats.averageRewardRate}% cashback automatically
                    </Typography>
                  </Box>
                )}

                <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                  <Chip 
                    label={bill.status.charAt(0).toUpperCase() + bill.status.slice(1)} 
                    color={getStatusColor(bill.status) as any}
                    size="small" 
                  />
                  {bill.untetherEnabled && (
                    <Chip 
                      label="Untether Enabled" 
                      color="success" 
                      size="small" 
                      icon={<AutoAwesome />}
                    />
                  )}
                </Box>

                <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
                  <Box sx={{ display: 'flex', alignItems: 'center', gap: 0.5 }}>
                    {getStatusIcon(bill.status)}
                    <Typography variant="body2" color="text.secondary">
                      {bill.status === 'upcoming' 
                        ? `${getDaysUntilDue(bill.dueDate)} days until due`
                        : bill.status === 'paid'
                        ? 'Paid automatically'
                        : 'Overdue'
                      }
                    </Typography>
                  </Box>
                </Box>

                {bill.lastPaid && (
                  <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                    Last paid: {formatDate(bill.lastPaid)}
                  </Typography>
                )}
              </CardContent>
              
              <CardActions>
                <Button 
                  size="small" 
                  variant="contained" 
                  sx={{ bgcolor: '#1B4D3E' }}
                  onClick={() => handlePayBill(bill.id)}
                >
                  {bill.untetherEnabled ? 'Pay with Untether' : 'Pay Now'}
                </Button>
                <Button size="small" color="primary">
                  Edit
                </Button>
                <Button size="small" color="secondary">
                  History
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* How It Works Section */}
      <Paper sx={{ mt: 4, p: 3, bgcolor: '#f8f9fa' }}>
        <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
          How Untether Credit Card Rewards Work
        </Typography>
        
        <Grid container spacing={3}>
          <Grid item xs={12} md={4}>
            <Box sx={{ textAlign: 'center' }}>
              <Avatar sx={{ bgcolor: '#1B4D3E', width: 60, height: 60, mx: 'auto', mb: 2 }}>
                <CreditCard sx={{ fontSize: 30 }} />
              </Avatar>
              <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                1. Enable Untether
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Turn on Untether for any bill to start earning rewards automatically
              </Typography>
            </Box>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Box sx={{ textAlign: 'center' }}>
              <Avatar sx={{ bgcolor: '#1B4D3E', width: 60, height: 60, mx: 'auto', mb: 2 }}>
                <AutoAwesome sx={{ fontSize: 30 }} />
              </Avatar>
              <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                2. Automatic Payment
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Untether pays your bills with your credit card and credits your checking account
              </Typography>
            </Box>
          </Grid>
          
          <Grid item xs={12} md={4}>
            <Box sx={{ textAlign: 'center' }}>
              <Avatar sx={{ bgcolor: '#1B4D3E', width: 60, height: 60, mx: 'auto', mb: 2 }}>
                <Star sx={{ fontSize: 30 }} />
              </Avatar>
              <Typography variant="subtitle1" sx={{ fontWeight: 600, mb: 1 }}>
                3. Earn Rewards
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Get {untetherStats.averageRewardRate}% cashback on all bills paid through Untether
              </Typography>
            </Box>
          </Grid>
        </Grid>
      </Paper>

      {/* Add Bill Dialog */}
      <Dialog open={addBillOpen} onClose={handleCloseAddBill} maxWidth="sm" fullWidth>
        <DialogTitle>Add New Bill</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Bill Name"
              placeholder="e.g., Rent Payment"
              fullWidth
            />
            <TextField
              label="Amount"
              type="number"
              placeholder="0.00"
              fullWidth
            />
            <TextField
              label="Due Date"
              type="date"
              InputLabelProps={{ shrink: true }}
              fullWidth
            />
            <TextField
              label="Category"
              select
              fullWidth
              SelectProps={{ native: true }}
            >
              <option value="">Select Category</option>
              <option value="housing">Housing</option>
              <option value="utilities">Utilities</option>
              <option value="insurance">Insurance</option>
              <option value="loans">Loans</option>
              <option value="subscriptions">Subscriptions</option>
              <option value="other">Other</option>
            </TextField>
            <FormControlLabel
              control={<Switch defaultChecked color="success" />}
              label="Enable Untether Credit Card Rewards"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseAddBill}>Cancel</Button>
          <Button variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
            Add Bill
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}