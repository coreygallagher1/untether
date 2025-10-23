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
  Alert,
  LinearProgress,
  Avatar,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  TextField,
  Switch,
  FormControlLabel,
  Tabs,
  Tab,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TablePagination
} from '@mui/material';
import { 
  Receipt, 
  Upload,
  CheckCircle,
  Warning,
  Info,
  FileDownload,
  FileUpload,
  Folder,
  Share,
  PersonAdd,
  Security,
  CloudUpload,
  Visibility,
  Edit,
  Delete,
  Search,
  FilterList,
  Favorite,
  AttachMoney,
  CalendarToday
} from '@mui/icons-material';
import { useState } from 'react';

export default function TaxCenterPage() {
  const [selectedYear, setSelectedYear] = useState('2024');
  const [activeTab, setActiveTab] = useState(0);
  const [uploadDialogOpen, setUploadDialogOpen] = useState(false);
  const [shareDialogOpen, setShareDialogOpen] = useState(false);
  const [donationDialogOpen, setDonationDialogOpen] = useState(false);
  const [page, setPage] = useState(0);
  const [rowsPerPage, setRowsPerPage] = useState(10);

  // Sample tax data
  const taxSummary = {
    totalDocuments: 15,
    documentsUploaded: 12,
    documentsPending: 3,
    totalSize: '1.8 GB',
    lastUpdated: '2024-12-15',
    accountantAccess: true,
    accountantName: 'Sarah Johnson CPA',
    totalDonations: 2450.00,
    donationCount: 8
  };

  const taxDocuments = [
    { 
      id: 1,
      name: 'Mortgage Interest Statement', 
      year: '2024',
      type: 'Deduction',
      category: 'Mortgage',
      status: 'uploaded',
      size: '156 KB',
      uploadDate: '2024-01-31',
      description: 'Annual mortgage interest paid',
      tags: ['deduction', 'mortgage', 'home']
    },
    { 
      id: 2,
      name: 'Property Tax Receipt', 
      year: '2024',
      type: 'Deduction',
      category: 'Property Tax',
      status: 'uploaded',
      size: '98 KB',
      uploadDate: '2024-01-31',
      description: 'Annual property tax payment',
      tags: ['deduction', 'property', 'tax']
    },
    { 
      id: 3,
      name: 'Charitable Donations Summary', 
      year: '2024',
      type: 'Deduction',
      category: 'Charitable',
      status: 'uploaded',
      size: '423 KB',
      uploadDate: '2024-01-31',
      description: 'Receipts for charitable contributions',
      tags: ['deduction', 'charitable', 'donations']
    },
    { 
      id: 4,
      name: 'Medical Expenses', 
      year: '2024',
      type: 'Deduction',
      category: 'Medical',
      status: 'uploaded',
      size: '567 KB',
      uploadDate: '2024-01-31',
      description: 'Medical bills and prescription receipts',
      tags: ['deduction', 'medical', 'health']
    },
    { 
      id: 5,
      name: 'Business Expenses', 
      year: '2024',
      type: 'Deduction',
      category: 'Business',
      status: 'uploaded',
      size: '789 KB',
      uploadDate: '2024-01-31',
      description: 'Home office and business-related expenses',
      tags: ['deduction', 'business', 'office']
    },
    { 
      id: 6,
      name: 'Student Loan Interest', 
      year: '2024',
      type: 'Deduction',
      category: 'Student Loan',
      status: 'pending',
      size: '0 KB',
      uploadDate: null,
      description: 'Student loan interest statement',
      tags: ['deduction', 'student', 'loan']
    },
    { 
      id: 7,
      name: 'HSA Contributions', 
      year: '2024',
      type: 'Deduction',
      category: 'HSA',
      status: 'pending',
      size: '0 KB',
      uploadDate: null,
      description: 'Health Savings Account contributions',
      tags: ['deduction', 'hsa', 'health']
    },
    { 
      id: 8,
      name: 'Education Expenses', 
      year: '2024',
      type: 'Deduction',
      category: 'Education',
      status: 'pending',
      size: '0 KB',
      uploadDate: null,
      description: 'Tuition and education-related expenses',
      tags: ['deduction', 'education', 'tuition']
    }
  ];

  const donations = [
    {
      id: 1,
      organization: 'American Red Cross',
      amount: 500.00,
      date: '2024-03-15',
      type: 'Cash',
      description: 'Disaster relief fund',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 2,
      organization: 'Local Food Bank',
      amount: 250.00,
      date: '2024-06-20',
      type: 'Cash',
      description: 'Monthly donation',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 3,
      organization: 'St. Jude Children\'s Hospital',
      amount: 1000.00,
      date: '2024-09-10',
      type: 'Cash',
      description: 'Annual donation',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 4,
      organization: 'Habitat for Humanity',
      amount: 300.00,
      date: '2024-11-05',
      type: 'Cash',
      description: 'Building materials fund',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 5,
      organization: 'Local Animal Shelter',
      amount: 200.00,
      date: '2024-12-01',
      type: 'Cash',
      description: 'Pet care fund',
      receipt: 'pending',
      deductible: true
    },
    {
      id: 6,
      organization: 'Doctors Without Borders',
      amount: 150.00,
      date: '2024-08-15',
      type: 'Cash',
      description: 'Medical aid',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 7,
      organization: 'Local Library',
      amount: 50.00,
      date: '2024-07-22',
      type: 'Cash',
      description: 'Book fund',
      receipt: 'uploaded',
      deductible: true
    },
    {
      id: 8,
      organization: 'Environmental Defense Fund',
      amount: 100.00,
      date: '2024-10-30',
      type: 'Cash',
      description: 'Climate action',
      receipt: 'uploaded',
      deductible: true
    }
  ];

  const accountantAccess = [
    {
      id: 1,
      name: 'Sarah Johnson CPA',
      email: 'sarah.johnson@taxpro.com',
      accessLevel: 'Full Access',
      lastLogin: '2024-12-14',
      permissions: ['view', 'download', 'comment'],
      status: 'active'
    },
    {
      id: 2,
      name: 'Mike Chen CPA',
      email: 'mike.chen@taxpro.com',
      accessLevel: 'View Only',
      lastLogin: '2024-12-10',
      permissions: ['view'],
      status: 'active'
    }
  ];

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

  const handleAddDonation = () => {
    setDonationDialogOpen(true);
  };

  const handleCloseDonation = () => {
    setDonationDialogOpen(false);
  };

  const handleChangePage = (event: unknown, newPage: number) => {
    setPage(newPage);
  };

  const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
    setRowsPerPage(parseInt(event.target.value, 10));
    setPage(0);
  };

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'uploaded': return 'success';
      case 'pending': return 'warning';
      case 'overdue': return 'error';
      default: return 'default';
    }
  };

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'uploaded': return <CheckCircle sx={{ color: 'success.main' }} />;
      case 'pending': return <Warning sx={{ color: 'warning.main' }} />;
      case 'overdue': return <Warning sx={{ color: 'error.main' }} />;
      default: return <Info />;
    }
  };

  const getTypeColor = (type: string) => {
    switch (type) {
      case 'Income': return 'primary';
      case 'Deduction': return 'secondary';
      default: return 'default';
    }
  };

  const handleUploadDocument = () => {
    setUploadDialogOpen(true);
  };

  const handleShareAccess = () => {
    setShareDialogOpen(true);
  };

  const handleCloseUpload = () => {
    setUploadDialogOpen(false);
  };

  const handleCloseShare = () => {
    setShareDialogOpen(false);
  };

  const uploadedDocs = taxDocuments.filter(doc => doc.status === 'uploaded');
  const pendingDocs = taxDocuments.filter(doc => doc.status === 'pending');

  return (
    <Box sx={{ p: 3, maxWidth: 1400, mx: 'auto' }}>
      {/* Header */}
      <Box sx={{ mb: 4 }}>
        <Typography variant="h4" component="h1" sx={{ 
          fontWeight: 600, 
          color: '#1B4D3E',
          mb: 1
        }}>
          Tax Document Center
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Secure storage and organization for all your tax documents with accountant access
        </Typography>
      </Box>

      {/* Tax Year Selector */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2, alignItems: 'center' }}>
        <Typography variant="h6" sx={{ fontWeight: 600 }}>
          Tax Year:
        </Typography>
        <Chip 
          label="2024" 
          color="primary" 
          variant="outlined"
          onClick={() => setSelectedYear('2024')}
        />
        <Chip 
          label="2023" 
          color="default" 
          variant="outlined"
          onClick={() => setSelectedYear('2023')}
        />
        <Chip 
          label="2022" 
          color="default" 
          variant="outlined"
          onClick={() => setSelectedYear('2022')}
        />
      </Box>

      {/* Summary Cards */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Folder sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Total Documents
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {taxSummary.totalDocuments}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {taxSummary.totalSize} stored
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <CloudUpload sx={{ fontSize: 24, color: 'success.main', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Uploaded
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: 'success.main' }}>
                {taxSummary.documentsUploaded}
              </Typography>
              <LinearProgress 
                variant="determinate" 
                value={(taxSummary.documentsUploaded / taxSummary.totalDocuments) * 100} 
                sx={{ 
                  mt: 1,
                  height: 6, 
                  borderRadius: 3,
                  bgcolor: '#e0e0e0',
                  '& .MuiLinearProgress-bar': {
                    bgcolor: 'success.main'
                  }
                }} 
              />
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Favorite sx={{ fontSize: 24, color: '#e91e63', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Total Donations
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#e91e63' }}>
                {formatCurrency(taxSummary.totalDonations)}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                {taxSummary.donationCount} donations
              </Typography>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} sm={6} md={3}>
          <Card sx={{ height: '100%', border: '1px solid #e0e0e0' }}>
            <CardContent>
              <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                <Security sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
                <Typography variant="subtitle2" color="text.secondary">
                  Accountant Access
                </Typography>
              </Box>
              <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
                {accountantAccess.length}
              </Typography>
              <Typography variant="body2" color="text.secondary">
                Active users
              </Typography>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Quick Actions */}
      <Box sx={{ mb: 3, display: 'flex', gap: 2, flexWrap: 'wrap' }}>
        <Button
          variant="contained"
          startIcon={<CloudUpload />}
          onClick={handleUploadDocument}
          sx={{
            bgcolor: '#1B4D3E',
            '&:hover': {
              bgcolor: '#143C30'
            }
          }}
        >
          Upload Documents
        </Button>
        <Button
          variant="contained"
          startIcon={<Favorite />}
          onClick={handleAddDonation}
          sx={{
            bgcolor: '#e91e63',
            '&:hover': {
              bgcolor: '#c2185b'
            }
          }}
        >
          Add Donation
        </Button>
        <Button
          variant="outlined"
          startIcon={<PersonAdd />}
          onClick={handleShareAccess}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Share with Accountant
        </Button>
        <Button
          variant="outlined"
          startIcon={<Search />}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Search Documents
        </Button>
        <Button
          variant="outlined"
          startIcon={<FilterList />}
          sx={{
            color: '#1B4D3E',
            borderColor: '#1B4D3E',
            '&:hover': {
              borderColor: '#1B4D3E',
              bgcolor: 'rgba(27, 77, 62, 0.1)'
            }
          }}
        >
          Filter by Type
        </Button>
      </Box>

      {/* Main Content Tabs */}
      <Paper sx={{ mb: 3 }}>
        <Tabs 
          value={activeTab} 
          onChange={(e, newValue) => setActiveTab(newValue)}
          sx={{ borderBottom: 1, borderColor: 'divider' }}
        >
          <Tab 
            label={`All Documents (${taxSummary.totalDocuments})`} 
            icon={<Folder />}
            iconPosition="start"
          />
          <Tab 
            label={`Uploaded (${taxSummary.documentsUploaded})`} 
            icon={<CheckCircle />}
            iconPosition="start"
          />
          <Tab 
            label={`Pending (${taxSummary.documentsPending})`} 
            icon={<Warning />}
            iconPosition="start"
          />
          <Tab 
            label={`Donations (${taxSummary.donationCount})`} 
            icon={<Favorite />}
            iconPosition="start"
          />
          <Tab 
            label={`Accountant Access`} 
            icon={<Share />}
            iconPosition="start"
          />
        </Tabs>
      </Paper>

      {/* Documents Grid */}
      <Grid container spacing={3}>
        {activeTab === 0 && (
          <>
            {taxDocuments.map((doc) => (
              <Grid item xs={12} md={6} lg={4} key={doc.id}>
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
                        bgcolor: doc.status === 'uploaded' ? 'success.main' : 'warning.main', 
                        width: 40, 
                        height: 40,
                        mr: 2
                      }}>
                        <Receipt />
                      </Avatar>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                          {doc.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {doc.category} • {doc.year}
                        </Typography>
                      </Box>
                    </Box>
                    
                    <Typography variant="body2" sx={{ mb: 2 }}>
                      {doc.description}
                    </Typography>

                    <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
                      <Chip 
                        label={doc.type} 
                        size="small" 
                        color={getTypeColor(doc.type) as any}
                      />
                      <Chip 
                        label={doc.status} 
                        size="small" 
                        color={getStatusColor(doc.status) as any}
                      />
                      <Chip 
                        label={doc.size} 
                        size="small" 
                        variant="outlined"
                      />
                    </Box>

                    {doc.uploadDate && (
                      <Typography variant="body2" color="text.secondary">
                        Uploaded: {formatDate(doc.uploadDate)}
                      </Typography>
                    )}
                  </CardContent>
                  
                  <CardActions>
                    <Button size="small" color="primary" startIcon={<Visibility />}>
                      View
                    </Button>
                    <Button size="small" color="primary" startIcon={<FileDownload />}>
                      Download
                    </Button>
                    <Button size="small" color="secondary" startIcon={<Edit />}>
                      Edit
                    </Button>
                    <Button size="small" color="error" startIcon={<Delete />}>
                      Delete
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </>
        )}

        {activeTab === 1 && (
          <>
            {uploadedDocs.map((doc) => (
              <Grid item xs={12} md={6} lg={4} key={doc.id}>
                <Card sx={{ 
                  height: '100%',
                  border: '2px solid #4caf50',
                  '&:hover': {
                    boxShadow: '0 4px 12px rgba(76, 175, 80, 0.2)',
                    transform: 'translateY(-2px)',
                    transition: 'all 0.2s ease-in-out'
                  }
                }}>
                  <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Avatar sx={{ 
                        bgcolor: 'success.main', 
                        width: 40, 
                        height: 40,
                        mr: 2
                      }}>
                        <CheckCircle />
                      </Avatar>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                          {doc.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {doc.category} • {doc.year}
                        </Typography>
                      </Box>
                    </Box>
                    
                    <Typography variant="body2" sx={{ mb: 2 }}>
                      {doc.description}
                    </Typography>

                    <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
                      <Chip 
                        label={doc.type} 
                        size="small" 
                        color={getTypeColor(doc.type) as any}
                      />
                      <Chip 
                        label={doc.size} 
                        size="small" 
                        variant="outlined"
                      />
                    </Box>

                    <Typography variant="body2" color="text.secondary">
                      Uploaded: {formatDate(doc.uploadDate!)}
                    </Typography>
                  </CardContent>
                  
                  <CardActions>
                    <Button size="small" color="primary" startIcon={<Visibility />}>
                      View
                    </Button>
                    <Button size="small" color="primary" startIcon={<FileDownload />}>
                      Download
                    </Button>
                    <Button size="small" color="secondary" startIcon={<Edit />}>
                      Edit
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </>
        )}

        {activeTab === 2 && (
          <>
            {pendingDocs.map((doc) => (
              <Grid item xs={12} md={6} lg={4} key={doc.id}>
                <Card sx={{ 
                  height: '100%',
                  border: '2px solid #ff9800',
                  '&:hover': {
                    boxShadow: '0 4px 12px rgba(255, 152, 0, 0.2)',
                    transform: 'translateY(-2px)',
                    transition: 'all 0.2s ease-in-out'
                  }
                }}>
                  <CardContent>
                    <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
                      <Avatar sx={{ 
                        bgcolor: 'warning.main', 
                        width: 40, 
                        height: 40,
                        mr: 2
                      }}>
                        <Warning />
                      </Avatar>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                          {doc.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {doc.category} • {doc.year}
                        </Typography>
                      </Box>
                    </Box>
                    
                    <Typography variant="body2" sx={{ mb: 2 }}>
                      {doc.description}
                    </Typography>

                    <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
                      <Chip 
                        label={doc.type} 
                        size="small" 
                        color={getTypeColor(doc.type) as any}
                      />
                      <Chip 
                        label="Pending Upload" 
                        size="small" 
                        color="warning"
                      />
                    </Box>

                    <Alert severity="warning" sx={{ mt: 2 }}>
                      <Typography variant="body2">
                        This document is still needed for your {doc.year} tax filing.
                      </Typography>
                    </Alert>
                  </CardContent>
                  
                  <CardActions>
                    <Button size="small" variant="contained" startIcon={<CloudUpload />} sx={{ bgcolor: '#1B4D3E' }}>
                      Upload Now
                    </Button>
                    <Button size="small" color="secondary">
                      Mark as Not Applicable
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </>
        )}

        {activeTab === 3 && (
          <>
            {/* Donations Table */}
            <Grid item xs={12}>
              <Card sx={{ border: '1px solid #e0e0e0' }}>
                <CardContent>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center', mb: 2 }}>
                    <Typography variant="h6" sx={{ fontWeight: 600 }}>
                      Charitable Donations - {selectedYear}
                    </Typography>
                    <Button
                      variant="contained"
                      startIcon={<Favorite />}
                      onClick={handleAddDonation}
                      sx={{
                        bgcolor: '#e91e63',
                        '&:hover': {
                          bgcolor: '#c2185b'
                        }
                      }}
                    >
                      Add Donation
                    </Button>
                  </Box>
                  
                  <TableContainer>
                    <Table>
                      <TableHead>
                        <TableRow>
                          <TableCell><strong>Organization</strong></TableCell>
                          <TableCell><strong>Amount</strong></TableCell>
                          <TableCell><strong>Date</strong></TableCell>
                          <TableCell><strong>Type</strong></TableCell>
                          <TableCell><strong>Description</strong></TableCell>
                          <TableCell><strong>Receipt</strong></TableCell>
                          <TableCell><strong>Actions</strong></TableCell>
                        </TableRow>
                      </TableHead>
                      <TableBody>
                        {donations
                          .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                          .map((donation) => (
                          <TableRow key={donation.id}>
                            <TableCell>
                              <Box sx={{ display: 'flex', alignItems: 'center' }}>
                                <Avatar sx={{ 
                                  bgcolor: '#e91e63', 
                                  width: 32, 
                                  height: 32,
                                  mr: 2,
                                  fontSize: '0.875rem'
                                }}>
                                  <Favorite sx={{ fontSize: 16 }} />
                                </Avatar>
                                <Typography variant="body2" sx={{ fontWeight: 500 }}>
                                  {donation.organization}
                                </Typography>
                              </Box>
                            </TableCell>
                            <TableCell>
                              <Typography variant="body2" sx={{ fontWeight: 600, color: '#e91e63' }}>
                                {formatCurrency(donation.amount)}
                              </Typography>
                            </TableCell>
                            <TableCell>
                              <Typography variant="body2">
                                {formatDate(donation.date)}
                              </Typography>
                            </TableCell>
                            <TableCell>
                              <Chip 
                                label={donation.type} 
                                size="small" 
                                color="primary"
                                variant="outlined"
                              />
                            </TableCell>
                            <TableCell>
                              <Typography variant="body2" color="text.secondary">
                                {donation.description}
                              </Typography>
                            </TableCell>
                            <TableCell>
                              <Chip 
                                label={donation.receipt} 
                                size="small" 
                                color={donation.receipt === 'uploaded' ? 'success' : 'warning'}
                              />
                            </TableCell>
                            <TableCell>
                              <Box sx={{ display: 'flex', gap: 1 }}>
                                <Button size="small" color="primary" startIcon={<Visibility />}>
                                  View
                                </Button>
                                <Button size="small" color="secondary" startIcon={<Edit />}>
                                  Edit
                                </Button>
                                <Button size="small" color="error" startIcon={<Delete />}>
                                  Delete
                                </Button>
                              </Box>
                            </TableCell>
                          </TableRow>
                        ))}
                      </TableBody>
                    </Table>
                  </TableContainer>
                  
                  <TablePagination
                    rowsPerPageOptions={[5, 10, 25]}
                    component="div"
                    count={donations.length}
                    rowsPerPage={rowsPerPage}
                    page={page}
                    onPageChange={handleChangePage}
                    onRowsPerPageChange={handleChangeRowsPerPage}
                  />
                </CardContent>
              </Card>
            </Grid>
          </>
        )}

        {activeTab === 4 && (
          <>
            {accountantAccess.map((accountant) => (
              <Grid item xs={12} md={6} key={accountant.id}>
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
                        bgcolor: '#1B4D3E', 
                        width: 48, 
                        height: 48,
                        mr: 2
                      }}>
                        {accountant.name.split(' ').map(n => n[0]).join('')}
                      </Avatar>
                      <Box sx={{ flexGrow: 1 }}>
                        <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                          {accountant.name}
                        </Typography>
                        <Typography variant="body2" color="text.secondary">
                          {accountant.email}
                        </Typography>
                      </Box>
                      <Chip 
                        label={accountant.status} 
                        color="success" 
                        size="small"
                      />
                    </Box>
                    
                    <Box sx={{ mb: 2 }}>
                      <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                        Access Level: {accountant.accessLevel}
                      </Typography>
                      <Typography variant="body2" color="text.secondary">
                        Last Login: {formatDate(accountant.lastLogin)}
                      </Typography>
                    </Box>

                    <Box sx={{ display: 'flex', gap: 1, mb: 2, flexWrap: 'wrap' }}>
                      {accountant.permissions.map((permission) => (
                        <Chip 
                          key={permission}
                          label={permission} 
                          size="small" 
                          color="primary"
                          variant="outlined"
                        />
                      ))}
                    </Box>
                  </CardContent>
                  
                  <CardActions>
                    <Button size="small" color="primary" startIcon={<Edit />}>
                      Edit Access
                    </Button>
                    <Button size="small" color="secondary" startIcon={<Visibility />}>
                      View Activity
                    </Button>
                    <Button size="small" color="error" startIcon={<Delete />}>
                      Revoke Access
                    </Button>
                  </CardActions>
                </Card>
              </Grid>
            ))}
          </>
        )}
      </Grid>

      {/* Upload Document Dialog */}
      <Dialog open={uploadDialogOpen} onClose={handleCloseUpload} maxWidth="sm" fullWidth>
        <DialogTitle>Upload Tax Document</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Document Name"
              placeholder="e.g., W-2 Form - Employer A"
              fullWidth
            />
            <TextField
              label="Document Type"
              select
              fullWidth
              SelectProps={{ native: true }}
            >
              <option value="">Select Type</option>
              <option value="w2">W-2 Form</option>
              <option value="1099-int">1099-INT</option>
              <option value="1099-div">1099-DIV</option>
              <option value="mortgage">Mortgage Interest</option>
              <option value="property">Property Tax</option>
              <option value="charitable">Charitable Donations</option>
              <option value="medical">Medical Expenses</option>
              <option value="business">Business Expenses</option>
              <option value="other">Other</option>
            </TextField>
            <TextField
              label="Description"
              multiline
              rows={3}
              placeholder="Brief description of the document..."
              fullWidth
            />
            <Box sx={{ border: '2px dashed #ccc', p: 3, textAlign: 'center', borderRadius: 1 }}>
              <CloudUpload sx={{ fontSize: 48, color: '#ccc', mb: 1 }} />
              <Typography variant="body1" sx={{ mb: 1 }}>
                Drag and drop files here, or click to browse
              </Typography>
              <Button variant="outlined" startIcon={<Upload />}>
                Choose Files
              </Button>
            </Box>
            <FormControlLabel
              control={<Switch defaultChecked />}
              label="Share with accountant automatically"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseUpload}>Cancel</Button>
          <Button variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
            Upload Document
          </Button>
        </DialogActions>
      </Dialog>

      {/* Share Access Dialog */}
      <Dialog open={shareDialogOpen} onClose={handleCloseShare} maxWidth="sm" fullWidth>
        <DialogTitle>Share Access with Accountant</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Accountant Name"
              placeholder="e.g., Sarah Johnson CPA"
              fullWidth
            />
            <TextField
              label="Email Address"
              type="email"
              placeholder="sarah.johnson@taxpro.com"
              fullWidth
            />
            <TextField
              label="Access Level"
              select
              fullWidth
              SelectProps={{ native: true }}
            >
              <option value="view">View Only</option>
              <option value="download">View & Download</option>
              <option value="full">Full Access</option>
            </TextField>
            <TextField
              label="Message (Optional)"
              multiline
              rows={3}
              placeholder="Add a personal message for your accountant..."
              fullWidth
            />
            <FormControlLabel
              control={<Switch defaultChecked />}
              label="Send email notification"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseShare}>Cancel</Button>
          <Button variant="contained" sx={{ bgcolor: '#1B4D3E' }}>
            Send Invitation
          </Button>
        </DialogActions>
      </Dialog>

      {/* Add Donation Dialog */}
      <Dialog open={donationDialogOpen} onClose={handleCloseDonation} maxWidth="sm" fullWidth>
        <DialogTitle>Add Charitable Donation</DialogTitle>
        <DialogContent>
          <Box sx={{ display: 'flex', flexDirection: 'column', gap: 2, pt: 1 }}>
            <TextField
              label="Organization Name"
              placeholder="e.g., American Red Cross"
              fullWidth
            />
            <TextField
              label="Donation Amount"
              type="number"
              placeholder="0.00"
              fullWidth
              InputProps={{
                startAdornment: <Typography sx={{ mr: 1 }}>$</Typography>
              }}
            />
            <TextField
              label="Donation Date"
              type="date"
              fullWidth
              InputLabelProps={{ shrink: true }}
            />
            <TextField
              label="Donation Type"
              select
              fullWidth
              SelectProps={{ native: true }}
            >
              <option value="cash">Cash</option>
              <option value="check">Check</option>
              <option value="credit">Credit Card</option>
              <option value="goods">Goods/Items</option>
              <option value="other">Other</option>
            </TextField>
            <TextField
              label="Description"
              multiline
              rows={3}
              placeholder="Brief description of the donation..."
              fullWidth
            />
            <Box sx={{ border: '2px dashed #ccc', p: 3, textAlign: 'center', borderRadius: 1 }}>
              <Receipt sx={{ fontSize: 48, color: '#ccc', mb: 1 }} />
              <Typography variant="body1" sx={{ mb: 1 }}>
                Upload receipt or acknowledgment letter
              </Typography>
              <Button variant="outlined" startIcon={<Upload />}>
                Choose File
              </Button>
            </Box>
            <FormControlLabel
              control={<Switch defaultChecked />}
              label="Tax deductible donation"
            />
          </Box>
        </DialogContent>
        <DialogActions>
          <Button onClick={handleCloseDonation}>Cancel</Button>
          <Button variant="contained" sx={{ bgcolor: '#e91e63' }}>
            Add Donation
          </Button>
        </DialogActions>
      </Dialog>
    </Box>
  );
}