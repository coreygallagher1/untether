'use client';

import { AppBar, Toolbar, Typography, Button, Box, IconButton, Menu, MenuItem, Avatar, Tabs, Tab } from '@mui/material';
import { Settings, ExitToApp, AccountBalance, Favorite, CreditCard, School, Receipt, Calculate } from '@mui/icons-material';
import { useState, useEffect } from 'react';
import { useRouter, usePathname } from 'next/navigation';
import Link from 'next/link';
import Image from 'next/image';
import { useAuth } from '../contexts/AuthContext';
import { getInitials } from '../../utils/initials';

export default function Navbar() {
  const router = useRouter();
  const pathname = usePathname();
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [activeTab, setActiveTab] = useState(0);
  const { isLoggedIn, logout, user } = useAuth();

  // Sync active tab with current pathname
  useEffect(() => {
    switch (pathname) {
      case '/dashboard':
        setActiveTab(-1); // Special value for dashboard (logo)
        break;
      case '/bank-accounts':
        setActiveTab(0);
        break;
      case '/bills':
        setActiveTab(1);
        break;
      case '/loans':
        setActiveTab(2);
        break;
      case '/causes':
        setActiveTab(3);
        break;
      case '/learn':
        setActiveTab(4);
        break;
      case '/tax-center':
        setActiveTab(5);
        break;
      default:
        // If on any other page, default to dashboard (logo)
        setActiveTab(-1);
        break;
    }
  }, [pathname]);

  const handleMenu = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const handleLogout = () => {
    logout();
    handleClose();
  };

  const handleSettings = () => {
    router.push('/settings');
    handleClose();
  };

  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setActiveTab(newValue);
    // Navigate to the appropriate page
    switch (newValue) {
      case 0:
        router.push('/bank-accounts');
        break;
      case 1:
        router.push('/bills');
        break;
      case 2:
        router.push('/loans');
        break;
      case 3:
        router.push('/causes');
        break;
      case 4:
        router.push('/learn');
        break;
      case 5:
        router.push('/tax-center');
        break;
      default:
        break;
    }
  };

  return (
    <AppBar 
      position="static" 
      sx={{ 
        bgcolor: '#1B4D3E',
        minHeight: 80,
        boxShadow: '0 2px 4px rgba(0,0,0,0.1)'
      }}
    >
      <Toolbar sx={{ minHeight: 80, py: 2 }}>
        <Box sx={{ 
          display: 'flex', 
          alignItems: 'center', 
          gap: 2,
          flexGrow: 1,
          '& img': {
            filter: 'drop-shadow(0 0 5px rgba(255, 255, 255, 0.6))',
            transition: 'filter 0.3s ease',
            '&:hover': {
              filter: 'drop-shadow(0 0 25px rgba(255, 255, 255, 0.8))',
            }
          }
        }}>
          <Box component={Link} href={isLoggedIn ? '/dashboard' : '/'} sx={{ 
            display: 'flex', 
            alignItems: 'center', 
            gap: '0.5rem',
            textDecoration: 'none',
            transition: 'all 0.3s ease',
            borderRadius: '8px',
            px: 2,
            py: 1,
            filter: activeTab === -1 ? 'drop-shadow(0 0 12px rgba(255, 255, 255, 0.4))' : 'none',
            '&:hover': {
              filter: 'drop-shadow(0 0 8px rgba(255, 255, 255, 0.3))'
            }
          }}>
            <Image
              src="/assets/UntetherLogo.png"
              alt="Untether Logo"
              width={56}
              height={56}
              style={{ objectFit: 'contain' }}
            />
            <Typography 
              variant="h6" 
              component="div" 
              sx={{ 
                fontSize: '1.75rem', 
                fontWeight: 600,
                fontFamily: 'var(--font-space-grotesk)',
                letterSpacing: '0.5px',
                color: 'white'
              }}
            >
              Untether
            </Typography>
          </Box>
          
          {/* Navigation Tabs - Only show when logged in */}
          {isLoggedIn && (
            <Tabs
              value={activeTab}
              onChange={handleTabChange}
              sx={{
                ml: 4,
                '& .MuiTab-root': {
                  color: 'rgba(255, 255, 255, 0.7)',
                  fontWeight: 500,
                  fontSize: '0.95rem',
                  minHeight: 48,
                  '&.Mui-selected': {
                    color: 'white',
                    fontWeight: 600,
                  },
                  '&:hover': {
                    color: 'rgba(255, 255, 255, 0.9)',
                  }
                },
                '& .MuiTabs-indicator': {
                  backgroundColor: 'white',
                  height: 3,
                  borderRadius: '2px 2px 0 0',
                }
              }}
            >
              <Tab 
                label="Bank Accounts" 
                icon={<AccountBalance sx={{ fontSize: 20 }} />}
                iconPosition="start"
              />
              <Tab 
                label="Bills" 
                icon={<Receipt sx={{ fontSize: 20 }} />}
                iconPosition="start"
              />
              <Tab 
                label="Loans" 
                icon={<CreditCard sx={{ fontSize: 20 }} />}
                iconPosition="start"
              />
              <Tab 
                label="Causes" 
                icon={<Favorite sx={{ fontSize: 20 }} />}
                iconPosition="start"
              />
                      <Tab 
                        label="Learn" 
                        icon={<School sx={{ fontSize: 20 }} />}
                        iconPosition="start"
                      />
                      <Tab 
                        label="Tax Center" 
                        icon={<Calculate sx={{ fontSize: 20 }} />}
                        iconPosition="start"
                      />
                    </Tabs>
          )}
        </Box>
        {isLoggedIn ? (
          <Box>
            <IconButton
              aria-label="account of current user"
              aria-controls="menu-appbar"
              aria-haspopup="true"
              onClick={handleMenu}
              sx={{ 
                p: 0.5,
                '&:hover': { 
                  bgcolor: 'rgba(255,255,255,0.1)' 
                },
                border: '2px solid rgba(255,255,255,0.2)',
                borderRadius: '50%'
              }}
            >
              <Avatar 
                sx={{ 
                  width: 44,
                  height: 44,
                  bgcolor: 'white',
                  color: '#1B4D3E',
                  fontSize: '1.25rem',
                  fontWeight: 600
                }}
              >
                {getInitials(user || {})}
              </Avatar>
            </IconButton>
            <Menu
              id="menu-appbar"
              anchorEl={anchorEl}
              anchorOrigin={{
                vertical: 'bottom',
                horizontal: 'right',
              }}
              keepMounted
              transformOrigin={{
                vertical: 'top',
                horizontal: 'right',
              }}
              open={Boolean(anchorEl)}
              onClose={handleClose}
              PaperProps={{
                sx: {
                  mt: 1,
                  '& .MuiMenuItem-root': {
                    py: 1,
                    px: 2,
                  },
                },
              }}
            >
                    <MenuItem onClick={handleSettings}>
                      <Settings sx={{ mr: 1 }} /> Settings
                    </MenuItem>
                    <MenuItem onClick={handleLogout}>
                      <ExitToApp sx={{ mr: 1 }} /> Logout
                    </MenuItem>
            </Menu>
          </Box>
        ) : (
          <Box sx={{ display: 'flex', gap: 2 }}>
            <Button
              component={Link}
              href="/auth/login"
              variant="outlined"
              sx={{
                color: 'white',
                borderColor: 'white',
                '&:hover': {
                  borderColor: 'white',
                  bgcolor: 'rgba(255, 255, 255, 0.1)'
                }
              }}
            >
              Log In
            </Button>
            <Button
              component={Link}
              href="/auth/signup"
              variant="contained"
              sx={{
                bgcolor: 'white',
                color: '#1B4D3E',
                '&:hover': {
                  bgcolor: 'rgba(255, 255, 255, 0.9)'
                }
              }}
            >
              Sign Up
            </Button>
          </Box>
        )}
      </Toolbar>
    </AppBar>
  );
} 