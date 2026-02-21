import { create } from 'zustand';
import type { User } from '../types';
import { authApi } from '../api/auth';
import { setTokens, clearTokens } from '../api/client';

interface AuthState {
  user: User | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  isInitialized: boolean;
  error: string | null;

  login: (login: string, password: string) => Promise<void>;
  register: (login: string, password: string) => Promise<void>;
  logout: () => void;
  fetchProfile: () => Promise<void>;
  init: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set) => ({
  user: null,
  isAuthenticated: false,
  isLoading: false,
  isInitialized: false,
  error: null,

  login: async (login, password) => {
    set({ isLoading: true, error: null });
    try {
      const res = await authApi.login(login, password);
      setTokens(res.access_token, res.refresh_token);
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (e) {
      set({ error: (e as Error).message, isLoading: false });
      throw e;
    }
  },

  register: async (login, password) => {
    set({ isLoading: true, error: null });
    try {
      const res = await authApi.register(login, password);
      setTokens(res.access_token, res.refresh_token);
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true, isLoading: false });
    } catch (e) {
      set({ error: (e as Error).message, isLoading: false });
      throw e;
    }
  },

  logout: () => {
    clearTokens();
    set({ user: null, isAuthenticated: false });
  },

  fetchProfile: async () => {
    try {
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true });
    } catch {
      clearTokens();
      set({ user: null, isAuthenticated: false });
    }
  },

  init: async () => {
    const token = localStorage.getItem('access_token');
    if (!token) {
      set({ isInitialized: true });
      return;
    }
    try {
      const user = await authApi.getProfile();
      set({ user, isAuthenticated: true, isInitialized: true });
    } catch {
      clearTokens();
      set({ user: null, isAuthenticated: false, isInitialized: true });
    }
  },
}));
