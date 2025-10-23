'use client';

import { 
  Box, 
  Typography, 
  Paper, 
  Grid, 
  Card, 
  CardContent, 
  Button, 
  Chip, 
  LinearProgress, 
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider,
  Avatar,
  IconButton
} from '@mui/material';
import { 
  AccountBalance, 
  TrendingUp, 
  TrendingDown,
  CreditCard, 
  Savings,
  Favorite,
  School,
  Notifications,
  Add,
  ArrowUpward,
  ArrowDownward,
  CheckCircle,
  Warning,
  Info
} from '@mui/icons-material';
import { useAuth } from '../contexts/AuthContext';
import { getInitials } from '../../utils/initials';

export default function DashboardPage() {
  const { user } = useAuth();

  // Sample data - in a real app, this would come from API calls
  const financialOverview = {
    totalBalance: 11268.99,
    monthlyIncome: 4500.00,
    monthlyExpenses: 3200.00,
    monthlySavings: 1300.00,
    roundupThisMonth: 47.23,
    creditScore: 742,
    emergencyFundProgress: 75,
    debtTotal: 31094.99
  };

  const recentTransactions = [
    { id: 1, description: 'Coffee Shop', amount: -4.50, date: 'Today', category: 'Food' },
    { id: 2, description: 'Roundup Donation', amount: 0.50, date: 'Today', category: 'Donation' },
    { id: 3, description: 'Salary Deposit', amount: 4500.00, date: 'Yesterday', category: 'Income' },
    { id: 4, description: 'Grocery Store', amount: -87.32, date: 'Yesterday', category: 'Food' },
    { id: 5, description: 'Roundup Savings', amount: 1.32, date: 'Yesterday', category: 'Savings' }
  ];

  const goals = [
    { title: 'Emergency Fund', target: 10000, current: 7500, deadline: 'Dec 2024', priority: 'high' },
    { title: 'Vacation Fund', target: 3000, current: 1200, deadline: 'Jun 2025', priority: 'medium' },
    { title: 'Debt Payoff', target: 31095, current: 18900, deadline: 'Dec 2026', priority: 'high' }
  ];

  const insights = [
    { type: 'success', message: 'Great job! You saved 15% more this month than last month.' },
    { type: 'warning', message: 'You\'re spending 20% more on dining out than your budget allows.' },
    { type: 'info', message: 'Your credit score improved by 12 points this quarter!' }
  ];

  const quickActions = [
    { title: 'Add Money', icon: <Add />, color: '#1B4D3E' },
    { title: 'Pay Bills', icon: <CreditCard />, color: '#1976d2' },
    { title: 'Transfer', icon: <TrendingUp />, color: '#388e3c' },
    { title: 'Invest', icon: <Savings />, color: '#f57c00' }
  ];

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD'
    }).format(amount);
  };

  const getInsightIcon = (type: string) => {
    switch (type) {
      case 'success': return <CheckCircle sx={{ color: 'success.main' }} />;
      case 'warning': return <Warning sx={{ color: 'warning.main' }} />;
      case 'info': return <Info sx={{ color: 'info.main' }} />;
      default: return <Info sx={{ color: 'info.main' }} />;
    }
  };

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'high': return 'error';
      case 'medium': return 'warning';
      case 'low': return 'success';
      default: return 'default';
    }
  };

  return (
    <Box sx={{ p: 3, maxWidth: 1400, mx: 'auto' }}>
      {/* Welcome Header */}
      <Paper sx={{ 
        p: 3, 
        mb: 4, 
        bgcolor: '#1B4D3E', 
        color: 'white',
        background: 'linear-gradient(135deg, #1B4D3E 0%, #2E7D32 100%)'
      }}>
        <Box sx={{ display: 'flex', alignItems: 'center', justifyContent: 'space-between' }}>
          <Box>
            <Typography variant="h4" component="h1" sx={{ fontWeight: 600, mb: 1 }}>
              Welcome back, {user?.first_name || 'User'}! 👋
            </Typography>
            <Typography variant="body1" sx={{ opacity: 0.9 }}>
              Here's your financial overview for today
            </Typography>
          </Box>
          <Box sx={{ display: 'flex', alignItems: 'center', gap: 2 }}>
            <IconButton sx={{ color: 'white' }}>
              <Notifications />
            </IconButton>
            <Avatar sx={{ bgcolor: 'white', color: '#1B4D3E' }}>
              {getInitials(user || {})}
            </Avatar>
          </Box>
        </Box>
      </Paper>

      {/* Financial Overview Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <AccountBalance sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Total Balance
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {formatCurrency(financialOverview.totalBalance)}
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                <ArrowUpward sx={{ fontSize: 16, color: 'success.main', mr: 0.5 }} />
                <Typography variant="body2" color="success.main">
                  +2.3% from last month
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <TrendingUp sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Monthly Income
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {formatCurrency(financialOverview.monthlyIncome)}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                Next deposit: Dec 1
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <TrendingDown sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Monthly Expenses
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {formatCurrency(financialOverview.monthlyExpenses)}
              </Typography>
              <Box sx={{ display: 'flex', alignItems: 'center', mt: 1 }}>
                <ArrowDownward sx={{ fontSize: 16, color: 'error.main', mr: 0.5 }} />
                <Typography variant="body2" color="error.main">
                  -5.2% from last month
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Savings sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Roundup This Month
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {formatCurrency(financialOverview.roundupThisMonth)}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mt: 1 }}>
                23 transactions
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      <Grid container spacing={3}>
        {/* Goals Progress */}
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
                Financial Goals
              </Typography>
              {goals.map((goal, index) => (
                <Box key={index} sx={{ mb: 3 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 1 }}>
                    <Typography variant="subtitle1" sx={{ fontWeight: 500 }}>
                      {goal.title}
                    </Typography>
                    <Chip 
                      label={goal.priority} 
                      size="small" 
                      color={getPriorityColor(goal.priority) as any}
                    />
                  </Box>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="body2" color="text.secondary">
                      {formatCurrency(goal.current)} / {formatCurrency(goal.target)}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {goal.deadline}
                    </Typography>
                  </Box>
                  <LinearProgress 
                    variant="determinate" 
                    value={(goal.current / goal.target) * 100} 
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
              ))}
              <Button variant="outlined" sx={{ color: '#1B4D3E', borderColor: '#1B4D3E' }}>
                View All Goals
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Recent Transactions */}
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
                Recent Activity
              </Typography>
              <List dense>
                {recentTransactions.map((transaction, index) => (
                  <ListItem key={transaction.id} sx={{ px: 0 }}>
                    <ListItemIcon sx={{ minWidth: 40 }}>
                      {transaction.amount > 0 ? (
                        <ArrowUpward sx={{ color: 'success.main' }} />
                      ) : (
                        <ArrowDownward sx={{ color: 'error.main' }} />
                      )}
                    </ListItemIcon>
                    <ListItemText
                      primary={transaction.description}
                      secondary={`${transaction.date} • ${transaction.category}`}
                    />
                    <Typography 
                      variant="body2" 
                      sx={{ 
                        fontWeight: 600,
                        color: transaction.amount > 0 ? 'success.main' : 'error.main'
                      }}
                    >
                      {transaction.amount > 0 ? '+' : ''}{formatCurrency(transaction.amount)}
                    </Typography>
                  </ListItem>
                ))}
              </List>
              <Divider sx={{ my: 2 }} />
              <Button variant="outlined" sx={{ color: '#1B4D3E', borderColor: '#1B4D3E' }}>
                View All Transactions
              </Button>
            </CardContent>
          </Card>
        </Grid>

        {/* Quick Actions */}
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
                Quick Actions
              </Typography>
              <Grid container spacing={2}>
                {quickActions.map((action, index) => (
                  <Grid item xs={6} key={index}>
                    <Button
                      variant="outlined"
                      startIcon={action.icon}
                      sx={{
                        width: '100%',
                        py: 2,
                        color: action.color,
                        borderColor: action.color,
                        '&:hover': {
                          borderColor: action.color,
                          bgcolor: `${action.color}10`
                        }
                      }}
                    >
                      {action.title}
                    </Button>
                  </Grid>
                ))}
              </Grid>
            </CardContent>
          </Card>
        </Grid>

        {/* Insights & Tips */}
        <Grid item xs={12} md={6}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Typography variant="h6" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
                Financial Insights
              </Typography>
              {insights.map((insight, index) => (
                <Box key={index} sx={{ display: 'flex', alignItems: 'flex-start', mb: 2 }}>
                  <Box sx={{ mr: 2, mt: 0.5 }}>
                    {getInsightIcon(insight.type)}
                  </Box>
                  <Typography variant="body2" sx={{ flex: 1 }}>
                    {insight.message}
                  </Typography>
                </Box>
              ))}
              <Button variant="outlined" sx={{ color: '#1B4D3E', borderColor: '#1B4D3E' }}>
                View All Insights
              </Button>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
}