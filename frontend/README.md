# Untether Frontend

This is the Next.js frontend application for Untether, a financial wellness platform that empowers individuals to take control of their financial future through automated savings, community support, and personalized guidance.

## Tech Stack

- **Framework:** Next.js 14 with App Router
- **Language:** TypeScript
- **UI Library:** Material-UI (MUI)
- **Styling:** CSS-in-JS with MUI's styled system
- **Form Management:** React Hook Form with Zod validation
- **State Management:** React Context API
- **HTTP Client:** Axios
- **Authentication:** JWT tokens

## Project Structure

```
frontend/
├── src/
│   ├── api/                   # API client
│   │   └── auth.ts           # Authentication API
│   ├── app/                  # App Router pages
│   │   ├── auth/            # Authentication pages
│   │   │   ├── login/       # Login page
│   │   │   └── signup/      # Signup page
│   │   ├── dashboard/       # Dashboard page
│   │   ├── settings/        # Settings page
│   │   ├── contexts/        # React contexts
│   │   │   └── AuthContext.tsx # Authentication context
│   │   ├── components/      # Shared components
│   │   │   ├── Favicon.tsx  # Dynamic favicon
│   │   │   ├── MainLayout.tsx # Main layout wrapper
│   │   │   └── Navbar.tsx   # Navigation component
│   │   └── styles/          # Global styles
│   │       └── theme.ts     # MUI theme configuration
│   └── components/          # Feature components
│       ├── features/        # Feature-specific components
│       │   └── auth/        # Authentication components
│       │       ├── LoginForm.tsx    # Login form
│       │       └── SignUpForm.tsx   # Signup form
│       └── layout/          # Layout components
├── public/                  # Static assets
│   └── assets/             # Images and icons
└── [config files]          # Next.js, TypeScript, ESLint configs
```

## Getting Started

### Prerequisites

- Node.js 18 or later
- npm, yarn, pnpm, or bun
- Backend services running (see backend README)

### Installation

1. Install dependencies:
   ```bash
   npm install
   ```

2. Start the development server:
   ```bash
   npm run dev
   ```

3. Open [http://localhost:3000](http://localhost:3000) in your browser

## Features

### Authentication
- **User Registration:** Complete signup flow with username, email, and password validation
- **User Login:** Secure login with email/username and password
- **Auto-login:** Automatic login after successful registration
- **Password Validation:** Real-time password strength feedback
- **JWT Integration:** Secure token-based authentication

### UI/UX
- **Material Design:** Modern, accessible UI components
- **Responsive Design:** Mobile-first responsive layout
- **Dark/Light Theme:** Dynamic theme switching
- **Form Validation:** Real-time validation with helpful error messages
- **Loading States:** Proper loading indicators and disabled states

### Development Features
- **TypeScript:** Full type safety
- **Hot Reload:** Instant updates during development
- **ESLint:** Code quality and consistency
- **Component Architecture:** Reusable, maintainable components

## API Integration

The frontend communicates with the Python FastAPI backend services:

- **User Service (Port 8001):** Authentication and user management
- **Plaid Service (Port 8002):** Bank account integration
- **Transaction Service (Port 8003):** Transaction processing

### API Client

The `src/api/auth.ts` file contains the API client for authentication:

```typescript
// Example usage
import { authApi } from '@/api/auth';

// Sign up
const response = await authApi.signUp({
  firstName: 'John',
  lastName: 'Doe',
  email: 'john@example.com',
  username: 'johndoe',
  password: 'SecurePassword123!'
});

// Login
const loginResponse = await authApi.login({
  email: 'john@example.com',
  password: 'SecurePassword123!'
});
```

## Available Scripts

- `npm run dev` - Start development server
- `npm run build` - Build for production
- `npm run start` - Start production server
- `npm run lint` - Run ESLint
- `npm run lint:fix` - Fix ESLint issues

## Environment Variables

Create a `.env.local` file for local development:

```env
NEXT_PUBLIC_API_URL=http://localhost:8001
```

## Learn More

- [Next.js Documentation](https://nextjs.org/docs)
- [Material-UI Documentation](https://mui.com/)
- [React Hook Form](https://react-hook-form.com/)
- [Zod Validation](https://zod.dev/)

## Contributing

1. Follow the existing code style and patterns
2. Use TypeScript for all new code
3. Write meaningful component and function names
4. Add proper error handling
5. Test your changes thoroughly