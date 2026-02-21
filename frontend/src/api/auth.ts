import type { AuthResponse, User } from '../types';
import { api } from './client';

export const authApi = {
  login: (login: string, password: string) =>
    api.post<AuthResponse>('/login', { login, password }),

  register: (login: string, password: string) =>
    api.post<AuthResponse>('/register', { login, password }),

  refreshToken: (refresh_token: string) =>
    api.post<{ access_token: string; refresh_token: string }>('/refresh-token', { refresh_token }),

  getProfile: () => api.get<User>('/profile'),
};
