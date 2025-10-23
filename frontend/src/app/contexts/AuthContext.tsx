import { createContext, useContext, useState, useEffect, ReactNode } from 'react';
import { useRouter } from 'next/navigation';
import { User } from '@/api/auth';
import { authApi } from '@/api/auth';

interface AuthContextType {
  isLoggedIn: boolean;
  user: User | null;
  token: string | null;
  login: (token: string, user: User) => void;
  logout: () => void;
  refreshUser: () => Promise<void>;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const router = useRouter();
  const [isLoggedIn, setIsLoggedIn] = useState(false);
  const [user, setUser] = useState<User | null>(null);
  const [token, setToken] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState(true);

  useEffect(() => {
    const initializeAuth = async () => {
      const storedToken = localStorage.getItem('token');
      if (storedToken) {
        try {
          console.log('Attempting to get user data with token:', storedToken.substring(0, 20) + '...');
          const userData = await authApi.getUser();
          console.log('Successfully got user data:', userData);
          setToken(storedToken);
          setUser(userData);
          setIsLoggedIn(true);
        } catch (error) {
          console.error('Failed to get user data:', error);
          // If we can't get user data, clear the invalid token and state
          localStorage.removeItem('token');
          setToken(null);
          setUser(null);
          setIsLoggedIn(false);
        }
      } else {
        // Ensure state is cleared if no token exists
        setToken(null);
        setUser(null);
        setIsLoggedIn(false);
      }
      setIsLoading(false);
    };

    initializeAuth();
  }, []);

  const login = (token: string, userData: User) => {
    localStorage.setItem('token', token);
    setToken(token);
    setUser(userData);
    setIsLoggedIn(true);
  };

  const logout = () => {
    localStorage.removeItem('token');
    setToken(null);
    setUser(null);
    setIsLoggedIn(false);
    router.push('/');
  };

  const refreshUser = async () => {
    try {
      const userData = await authApi.getUser();
      setUser(userData);
    } catch (error) {
      console.error('Failed to refresh user data:', error);
      // Don't logout automatically - just log the error
      // The user might still be logged in, just the refresh failed
    }
  };

  if (isLoading) {
    return null; // or return a loading spinner
  }

  return (
    <AuthContext.Provider value={{ isLoggedIn, user, token, login, logout, refreshUser }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
} 