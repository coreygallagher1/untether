import axios from 'axios';

const API_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';

interface SignUpData {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
}

interface LoginData {
  email: string;
  password: string;
}

export interface User {
  id: string;
  email: string;
  first_name: string;
  last_name: string;
  created_at: {
    seconds: number;
    nanos: number;
  };
  updated_at: {
    seconds: number;
    nanos: number;
  };
}

export interface AuthResponse {
  userId: string;
  token: string;
  user: User;
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
    this.baseUrl = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081';
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

      const response = await this.axiosInstance.get(`/api/v1/users/${userId}`);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data.message || 'Failed to get user data');
      }
      throw new Error('Failed to get user data');
    }
  }

  async signUp(data: SignUpData): Promise<AuthResponse> {
    try {
      const response = await this.axiosInstance.post('/api/v1/auth/signup', data);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data.message || 'Failed to sign up');
      }
      throw new Error('Failed to sign up');
    }
  }

  async login(data: LoginData): Promise<AuthResponse> {
    try {
      const response = await this.axiosInstance.post('/api/v1/auth/signin', data);
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data.message || 'Failed to login');
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

  async changePassword(userId: string, oldPassword: string, newPassword: string): Promise<ChangePasswordResponse> {
    try {
      const response = await this.axiosInstance.post('/api/v1/auth/change-password', {
        userId,
        oldPassword,
        newPassword,
      });
      return response.data;
    } catch (error) {
      if (axios.isAxiosError(error) && error.response) {
        throw new Error(error.response.data.message || 'Failed to change password');
      }
      throw new Error('Failed to change password');
    }
  }
}

export const authApi = new AuthApi(); 