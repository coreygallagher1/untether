import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8001';

interface SignUpData {
  firstName: string;
  lastName: string;
  email: string;
  username: string;
  password: string;
}

interface LoginData {
  email: string;
  password: string;
}

export interface User {
  id: number;
  email: string;
  username: string;
  first_name?: string;
  last_name?: string;
  is_active: boolean;
  is_verified: boolean;
  created_at: string;
}

export interface AuthResponse {
  access_token: string;
  token_type: string;
  user?: User; // For registration response
}

interface ResetPasswordResponse {
  success: boolean;
  message: string;
}

interface ChangePasswordResponse {
  success: boolean;
  message: string;
}

class AuthApi {
  private baseUrl: string;
  private axiosInstance;

  constructor() {
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8001';
    this.axiosInstance = axios.create({
      baseURL: this.baseUrl,
    });

    // Add token to requests if it exists
    this.axiosInstance.interceptors.request.use((config) => {
      const token = localStorage.getItem('token');
      if (token) {
        config.headers.Authorization = `Bearer ${token}`;
      }
      return config;
    });
  }

  async getUser(): Promise<User> {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        throw new Error('No token found');
      }

      // Decode the token to get the user ID
      const tokenPayload = JSON.parse(atob(token.split('.')[1]));
      const userId = tokenPayload.user_id;

      const response = await this.axiosInstance.get('/users/me');
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // Handle FastAPI error format: error.response.data.detail.message
        const errorMessage = error.response.data.detail?.message || error.response.data.message || 'Failed to get user data';
        throw new Error(errorMessage);
      }
      throw new Error('Failed to get user data');
    }
  }

  async signUp(data: SignUpData): Promise<AuthResponse> {
    try {
      const response = await this.axiosInstance.post('/auth/register', {
        email: data.email,
        username: data.username,
        first_name: data.firstName,
        last_name: data.lastName,
        password: data.password
      });
      // Registration now returns token directly, we need to get user data separately
      const token = response.data.access_token;
      const userResponse = await this.axiosInstance.get('/users/me', {
        headers: { Authorization: `Bearer ${token}` }
      });
      return {
        access_token: token,
        token_type: response.data.token_type,
        user: userResponse.data
      };
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // Handle FastAPI error format: error.response.data.detail.message
        const errorMessage = error.response.data.detail?.message || error.response.data.message || 'Failed to sign up';
        throw new Error(errorMessage);
      }
      throw new Error('Failed to sign up');
    }
  }

  async login(data: LoginData): Promise<AuthResponse> {
    try {
      const response = await this.axiosInstance.post('/auth/login', {
        username: data.email, // Use email as username for login
        password: data.password
      });
      // Login returns token, we need to get user data separately
      const token = response.data.access_token;
      const userResponse = await this.axiosInstance.get('/users/me', {
        headers: { Authorization: `Bearer ${token}` }
      });
      return {
        access_token: token,
        token_type: response.data.token_type,
        user: userResponse.data
      };
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        // Handle FastAPI error format: error.response.data.detail.message
        const errorMessage = error.response.data.detail?.message || error.response.data.message || 'Failed to login';
        throw new Error(errorMessage);
      }
      throw new Error('Failed to login');
    }
  }

  async resetPassword(email: string): Promise<ResetPasswordResponse> {
    try {
      const response = await this.axiosInstance.post('/api/v1/auth/reset-password', { email });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data.message || 'Failed to reset password');
      }
      throw new Error('Failed to reset password');
    }
  }

  async updateUser(userData: Partial<User>): Promise<User | { user: User; new_token: string; token_type: string }> {
    try {
      const response = await this.axiosInstance.put('/users/me', userData);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        console.error('API Error:', error.response.data);
        const errorMessage = error.response.data.detail?.message || error.response.data.message || 'Failed to update user';
        throw new Error(errorMessage);
      }
      throw new Error('Failed to update user');
    }
  }

  async changePassword(currentPassword: string, newPassword: string): Promise<{ message: string }> {
    try {
      const response = await this.axiosInstance.put('/users/me/password', {
        current_password: currentPassword,
        new_password: newPassword,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        const errorMessage = error.response.data.detail?.message || error.response.data.message || 'Failed to change password';
        throw new Error(errorMessage);
      }
      throw new Error('Failed to change password');
    }
  }

}

export const authApi = new AuthApi(); 