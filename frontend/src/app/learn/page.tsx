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
  LinearProgress, 
  Accordion, 
  AccordionSummary, 
  AccordionDetails,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Divider
} from '@mui/material';
import { 
  School, 
  ExpandMore, 
  CheckCircle, 
  PlayCircle, 
  Quiz, 
  TrendingUp, 
  AccountBalance, 
  CreditCard,
  Savings,
  Assessment
} from '@mui/icons-material';
import { useState } from 'react';

export default function LearnPage() {
  const [completedLessons, setCompletedLessons] = useState<number[]>([0, 1]); // Sample completed lessons
  const [currentLesson, setCurrentLesson] = useState(2); // Sample current lesson

  const courses = [
    {
      id: 1,
      title: "Budgeting Basics",
      description: "Learn how to create and maintain a personal budget",
      duration: "15 min",
      difficulty: "Beginner",
      lessons: [
        { title: "What is a Budget?", completed: true },
        { title: "50/30/20 Rule", completed: true },
        { title: "Tracking Your Expenses", completed: false },
        { title: "Budget Adjustments", completed: false }
      ],
      progress: 50,
      icon: <AccountBalance sx={{ fontSize: 32 }} />
    },
    {
      id: 2,
      title: "Understanding Credit",
      description: "Master the fundamentals of credit scores and credit cards",
      duration: "20 min",
      difficulty: "Beginner",
      lessons: [
        { title: "What is Credit?", completed: false },
        { title: "Credit Scores Explained", completed: false },
        { title: "Building Good Credit", completed: false },
        { title: "Credit Card Best Practices", completed: false }
      ],
      progress: 0,
      icon: <CreditCard sx={{ fontSize: 32 }} />
    },
    {
      id: 3,
      title: "Saving & Investing",
      description: "Build wealth through smart saving and investing strategies",
      duration: "25 min",
      difficulty: "Intermediate",
      lessons: [
        { title: "Emergency Fund", completed: false },
        { title: "Types of Savings Accounts", completed: false },
        { title: "Introduction to Investing", completed: false },
        { title: "Compound Interest", completed: false },
        { title: "Risk vs Return", completed: false }
      ],
      progress: 0,
      icon: <Savings sx={{ fontSize: 32 }} />
    },
    {
      id: 4,
      title: "Debt Management",
      description: "Strategies for paying off debt and staying debt-free",
      duration: "18 min",
      difficulty: "Intermediate",
      lessons: [
        { title: "Types of Debt", completed: false },
        { title: "Debt Snowball Method", completed: false },
        { title: "Debt Avalanche Method", completed: false },
        { title: "Avoiding New Debt", completed: false }
      ],
      progress: 0,
      icon: <TrendingUp sx={{ fontSize: 32 }} />
    }
  ];

  const quickTips = [
    {
      title: "The 24-Hour Rule",
      description: "Wait 24 hours before making any purchase over $100. This helps avoid impulse buying.",
      category: "Spending"
    },
    {
      title: "Pay Yourself First",
      description: "Set up automatic transfers to savings before paying other bills. Treat savings like a non-negotiable expense.",
      category: "Saving"
    },
    {
      title: "Emergency Fund Target",
      description: "Aim for 3-6 months of expenses in your emergency fund. Start with $1,000 if you're just beginning.",
      category: "Emergency"
    },
    {
      title: "Credit Utilization",
      description: "Keep your credit card usage below 30% of your available credit limit to maintain a good credit score.",
      category: "Credit"
    }
  ];

  const handleStartCourse = (courseId: number) => {
    // In a real app, this would navigate to the course content
    console.log(`Starting course ${courseId}`);
  };

  const handleResumeCourse = (courseId: number) => {
    // In a real app, this would navigate to the current lesson
    console.log(`Resuming course ${courseId}`);
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
          Financial Literacy Center
        </Typography>
        <Typography variant="body1" color="text.secondary">
          Build your financial knowledge with interactive courses and practical tips
        </Typography>
      </Box>

      {/* Progress Overview */}
      <Paper sx={{ p: 3, mb: 4, bgcolor: '#f8f9fa' }}>
        <Box sx={{ display: 'flex', alignItems: 'center', mb: 2 }}>
          <Assessment sx={{ fontSize: 24, color: '#1B4D3E', mr: 1 }} />
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            Your Learning Progress
          </Typography>
        </Box>
        
        <Grid container spacing={3}>
          <Grid item xs={12} sm={3}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              1
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Courses Completed
            </Typography>
          </Grid>
          <Grid item xs={12} sm={3}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              2
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Lessons Completed
            </Typography>
          </Grid>
          <Grid item xs={12} sm={3}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              15
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Minutes Learned
            </Typography>
          </Grid>
          <Grid item xs={12} sm={3}>
            <Typography variant="h5" sx={{ fontWeight: 600, color: '#1B4D3E' }}>
              Beginner
            </Typography>
            <Typography variant="body2" color="text.secondary">
              Current Level
            </Typography>
          </Grid>
        </Grid>
      </Paper>

      {/* Courses Section */}
      <Typography variant="h5" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
        Available Courses
      </Typography>

      <Grid container spacing={3} sx={{ mb: 4 }}>
        {courses.map((course) => (
          <Grid item xs={12} md={6} key={course.id}>
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
                  <Box sx={{ color: '#1B4D3E', mr: 2 }}>
                    {course.icon}
                  </Box>
                  <Box sx={{ flexGrow: 1 }}>
                    <Typography variant="h6" component="h2" sx={{ fontWeight: 600 }}>
                      {course.title}
                    </Typography>
                    <Typography variant="body2" color="text.secondary">
                      {course.description}
                    </Typography>
                  </Box>
                </Box>
                
                <Box sx={{ mb: 2 }}>
                  <Box sx={{ display: 'flex', justifyContent: 'space-between', mb: 1 }}>
                    <Typography variant="body2" color="text.secondary">
                      Progress
                    </Typography>
                    <Typography variant="body2" sx={{ fontWeight: 600 }}>
                      {course.progress}%
                    </Typography>
                  </Box>
                  <LinearProgress 
                    variant="determinate" 
                    value={course.progress} 
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
                    label={course.difficulty} 
                    color="primary" 
                    size="small" 
                  />
                  <Chip 
                    label={course.duration} 
                    color="secondary" 
                    size="small" 
                  />
                </Box>

                <Accordion sx={{ boxShadow: 'none', '&:before': { display: 'none' } }}>
                  <AccordionSummary expandIcon={<ExpandMore />}>
                    <Typography variant="body2" sx={{ fontWeight: 500 }}>
                      View Lessons ({course.lessons.length})
                    </Typography>
                  </AccordionSummary>
                  <AccordionDetails sx={{ pt: 0 }}>
                    <List dense>
                      {course.lessons.map((lesson, index) => (
                        <ListItem key={index} sx={{ py: 0.5 }}>
                          <ListItemIcon sx={{ minWidth: 32 }}>
                            {lesson.completed ? (
                              <CheckCircle sx={{ fontSize: 16, color: 'success.main' }} />
                            ) : (
                              <PlayCircle sx={{ fontSize: 16, color: 'text.secondary' }} />
                            )}
                          </ListItemIcon>
                          <ListItemText 
                            primary={lesson.title}
                            primaryTypographyProps={{
                              fontSize: '0.875rem',
                              color: lesson.completed ? 'success.main' : 'text.primary'
                            }}
                          />
                        </ListItem>
                      ))}
                    </List>
                  </AccordionDetails>
                </Accordion>
              </CardContent>
              
              <CardActions>
                {course.progress === 0 ? (
                  <Button 
                    size="small" 
                    variant="contained" 
                    sx={{ bgcolor: '#1B4D3E' }}
                    onClick={() => handleStartCourse(course.id)}
                  >
                    Start Course
                  </Button>
                ) : course.progress === 100 ? (
                  <Button size="small" color="success">
                    Completed
                  </Button>
                ) : (
                  <Button 
                    size="small" 
                    variant="contained" 
                    sx={{ bgcolor: '#1B4D3E' }}
                    onClick={() => handleResumeCourse(course.id)}
                  >
                    Continue
                  </Button>
                )}
                <Button size="small" color="primary">
                  Preview
                </Button>
              </CardActions>
            </Card>
          </Grid>
        ))}
      </Grid>

      {/* Quick Tips Section */}
      <Typography variant="h5" sx={{ fontWeight: 600, mb: 3, color: '#1B4D3E' }}>
        Quick Financial Tips
      </Typography>

      <Grid container spacing={2}>
        {quickTips.map((tip, index) => (
          <Grid item xs={12} sm={6} md={3} key={index}>
            <Paper sx={{ p: 2, height: '100%', border: '1px solid #e0e0e0' }}>
              <Typography variant="subtitle2" sx={{ fontWeight: 600, mb: 1, color: '#1B4D3E' }}>
                {tip.title}
              </Typography>
              <Typography variant="body2" color="text.secondary" sx={{ mb: 1 }}>
                {tip.description}
              </Typography>
              <Chip 
                label={tip.category} 
                size="small" 
                color="primary" 
                variant="outlined"
              />
            </Paper>
          </Grid>
        ))}
      </Grid>

      {/* Call to Action */}
      <Paper sx={{ mt: 4, p: 3, bgcolor: '#1B4D3E', color: 'white' }}>
        <Box sx={{ textAlign: 'center' }}>
          <School sx={{ fontSize: 48, mb: 2, opacity: 0.8 }} />
          <Typography variant="h6" sx={{ fontWeight: 600, mb: 1 }}>
            Ready to Improve Your Financial Future?
          </Typography>
          <Typography variant="body2" sx={{ mb: 2, opacity: 0.9 }}>
            Start with our Budgeting Basics course and build a solid foundation for financial success.
          </Typography>
          <Button 
            variant="contained" 
            sx={{ 
              bgcolor: 'white', 
              color: '#1B4D3E',
              '&:hover': {
                bgcolor: 'rgba(255, 255, 255, 0.9)'
              }
            }}
            onClick={() => handleStartCourse(1)}
          >
            Start Learning Now
          </Button>
        </Box>
      </Paper>
    </Box>
  );
}
