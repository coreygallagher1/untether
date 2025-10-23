'use client';

import { useAuth } from '@/app/contexts/AuthContext';
import { Box, Paper, Typography, Grid, Skeleton, Avatar, Button, Dialog, DialogTitle, DialogContent, DialogActions, TextField, Alert, Snackbar } from '@mui/material';
import { AccountCircle, Email, CalendarToday, Edit } from '@mui/icons-material';
import { getInitials } from '../../utils/initials';
import { useState } from 'react';
import { authApi } from '../../api/auth';

export default function SettingsPage() {
  const { user, refreshUser } = useAuth();
  
  // State for modals and forms
  const [editUsernameOpen, setEditUsernameOpen] = useState(false);
  const [editEmailOpen, setEditEmailOpen] = useState(false);
  const [changePasswordOpen, setChangePasswordOpen] = useState(false);
  
  // Form states
  const [newUsername, setNewUsername] = useState('');
  const [newEmail, setNewEmail] = useState('');
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  
  // Password validation criteria
  const passwordCriteria = [
    { label: 'At least 8 characters', met: newPassword.length >= 8 },
    { label: 'One uppercase letter', met: /[A-Z]/.test(newPassword) },
    { label: 'One lowercase letter', met: /[a-z]/.test(newPassword) },
    { label: 'One number', met: /[0-9]/.test(newPassword) },
    { label: 'One special character', met: /[^A-Za-z0-9]/.test(newPassword) },
  ];
  
  // Feedback states
  const [snackbarOpen, setSnackbarOpen] = useState(false);
  const [snackbarMessage, setSnackbarMessage] = useState('');
  const [snackbarSeverity, setSnackbarSeverity] = useState<'success' | 'error'>('success');
  const [loading, setLoading] = useState(false);

  console.log('Settings page - user data:', user);

  // Handler functions
  const showSnackbar = (message: string, severity: 'success' | 'error') => {
    setSnackbarMessage(message);
    setSnackbarSeverity(severity);
    setSnackbarOpen(true);
  };

  const handleEditUsername = () => {
    setNewUsername(user?.username || '');
    setEditUsernameOpen(true);
  };

  const handleEditEmail = () => {
    setNewEmail(user?.email || '');
    setEditEmailOpen(true);
  };

  const handleChangePassword = () => {
    setCurrentPassword('');
    setNewPassword('');
    setConfirmPassword('');
    setChangePasswordOpen(true);
  };

  const saveUsername = async () => {
    if (!newUsername.trim()) {
      showSnackbar('Username cannot be empty', 'error');
      return;
    }
    
    setLoading(true);
    try {
      const response = await authApi.updateUser({ username: newUsername });
      showSnackbar('Username updated successfully!', 'success');
      setEditUsernameOpen(false);
      
      // If we got a new token (username was updated), update it
      if ('new_token' in response) {
        localStorage.setItem('token', response.new_token);
        // Update the user data with the response
        await refreshUser();
      } else {
        // Just refresh user data normally
        await refreshUser();
      }
    } catch (error) {
      console.error('Username update error:', error);
      const errorMessage = error instanceof Error ? error.message : 'Failed to update username';
      showSnackbar(errorMessage, 'error');
      // Don't call refreshUser() on error - it might cause "User not found"
    } finally {
      setLoading(false);
    }
  };

  const saveEmail = async () => {
    if (!newEmail.trim()) {
      showSnackbar('Email cannot be empty', 'error');
      return;
    }
    
    setLoading(true);
    try {
      await authApi.updateUser({ email: newEmail });
      showSnackbar('Email updated successfully!', 'success');
      setEditEmailOpen(false);
      // Refresh user data
      await refreshUser();
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to update email';
      showSnackbar(errorMessage, 'error');
    } finally {
      setLoading(false);
    }
  };

  const savePassword = async () => {
    if (!currentPassword || !newPassword || !confirmPassword) {
      showSnackbar('All password fields are required', 'error');
      return;
    }
    
    if (newPassword !== confirmPassword) {
      showSnackbar('New passwords do not match', 'error');
      return;
    }
    
    // Check if password meets all criteria
    const allCriteriaMet = passwordCriteria.every(criterion => criterion.met);
    if (!allCriteriaMet) {
      showSnackbar('Password does not meet all requirements', 'error');
      return;
    }
    
    setLoading(true);
    try {
      await authApi.changePassword(currentPassword, newPassword);
      showSnackbar('Password changed successfully!', 'success');
      setChangePasswordOpen(false);
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (error) {
      const errorMessage = error instanceof Error ? error.message : 'Failed to change password';
      showSnackbar(errorMessage, 'error');
    } finally {
      setLoading(false);
    }
  };

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

  const formatDate = (dateString: string) => {
    try {
      const date = new Date(dateString);
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
            {getInitials(user)}
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
            <Grid item xs={12}>
              <Typography variant="subtitle2" color="text.secondary" gutterBottom sx={{ pl: 1 }}>
                Username
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
                  <AccountCircle sx={{ color: 'text.secondary', fontSize: 20 }} />
                  <Typography variant="body1" sx={{ fontWeight: 500 }}>
                    {user.username || 'Not available'}
                  </Typography>
                </Box>
                <Button 
                  startIcon={<Edit sx={{ fontSize: 18 }} />}
                  variant="outlined"
                  size="small"
                  onClick={handleEditUsername}
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
                  onClick={handleEditEmail}
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
              <Grid item xs={12}>
                <Button
                  fullWidth
                  variant="contained"
                  size="large"
                  onClick={handleChangePassword}
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

      {/* Edit Username Modal */}
      <Dialog open={editUsernameOpen} onClose={() => setEditUsernameOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Edit Username</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="New Username"
            fullWidth
            variant="outlined"
            value={newUsername}
            onChange={(e) => setNewUsername(e.target.value)}
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditUsernameOpen(false)}>Cancel</Button>
          <Button onClick={saveUsername} variant="contained" disabled={loading}>
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Edit Email Modal */}
      <Dialog open={editEmailOpen} onClose={() => setEditEmailOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Edit Email</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="New Email"
            type="email"
            fullWidth
            variant="outlined"
            value={newEmail}
            onChange={(e) => setNewEmail(e.target.value)}
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setEditEmailOpen(false)}>Cancel</Button>
          <Button onClick={saveEmail} variant="contained" disabled={loading}>
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Change Password Modal */}
      <Dialog open={changePasswordOpen} onClose={() => setChangePasswordOpen(false)} maxWidth="sm" fullWidth>
        <DialogTitle>Change Password</DialogTitle>
        <DialogContent>
          <TextField
            autoFocus
            margin="dense"
            label="Current Password"
            type="password"
            fullWidth
            variant="outlined"
            value={currentPassword}
            onChange={(e) => setCurrentPassword(e.target.value)}
            sx={{ mt: 2 }}
          />
          <TextField
            margin="dense"
            label="New Password"
            type="password"
            fullWidth
            variant="outlined"
            value={newPassword}
            onChange={(e) => setNewPassword(e.target.value)}
            sx={{ mt: 2 }}
          />
          
          {/* Password Criteria */}
          {newPassword && (
            <Box sx={{ mt: 2, p: 2, bgcolor: 'grey.50', borderRadius: 1 }}>
              <Typography variant="body2" sx={{ mb: 1, fontWeight: 'medium' }}>
                Password Requirements:
              </Typography>
              {passwordCriteria.map((criterion, index) => (
                <Box key={index} sx={{ display: 'flex', alignItems: 'center', mb: 0.5 }}>
                  <Box
                    sx={{
                      width: 16,
                      height: 16,
                      borderRadius: '50%',
                      bgcolor: criterion.met ? 'success.main' : 'grey.300',
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      mr: 1,
                    }}
                  >
                    {criterion.met && (
                      <Typography variant="caption" sx={{ color: 'white', fontSize: '10px' }}>
                        ✓
                      </Typography>
                    )}
                  </Box>
                  <Typography
                    variant="body2"
                    sx={{
                      color: criterion.met ? 'success.main' : 'text.secondary',
                      fontSize: '0.875rem',
                    }}
                  >
                    {criterion.label}
                  </Typography>
                </Box>
              ))}
            </Box>
          )}
          
          <TextField
            margin="dense"
            label="Confirm New Password"
            type="password"
            fullWidth
            variant="outlined"
            value={confirmPassword}
            onChange={(e) => setConfirmPassword(e.target.value)}
            sx={{ mt: 2 }}
          />
        </DialogContent>
        <DialogActions>
          <Button onClick={() => setChangePasswordOpen(false)}>Cancel</Button>
          <Button onClick={savePassword} variant="contained" disabled={loading}>
            {loading ? 'Saving...' : 'Save'}
          </Button>
        </DialogActions>
      </Dialog>

      {/* Success/Error Snackbar */}
      <Snackbar
        open={snackbarOpen}
        autoHideDuration={6000}
        onClose={() => setSnackbarOpen(false)}
      >
        <Alert onClose={() => setSnackbarOpen(false)} severity={snackbarSeverity}>
          {snackbarMessage}
        </Alert>
      </Snackbar>
    </Box>
  );
} 