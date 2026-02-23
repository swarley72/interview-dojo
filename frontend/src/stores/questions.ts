import { create } from 'zustand';
import type { QuestionShort, Difficulty, QuestionType } from '../types';
import { questionsApi } from '../api/questions';

interface QuestionsState {
  items: QuestionShort[];
  totalCount: number;
  page: number;
  limit: number;
  filterDifficulty: Difficulty | '';
  filterType: QuestionType | '';
  filterTagIds: number[];
  filterVerified: '' | 'true' | 'false';
  searchQuery: string;
  isLoading: boolean;
  isFetching: boolean;
  initialized: boolean;

  setFilter: (key: 'filterDifficulty' | 'filterType' | 'filterVerified', value: string) => void;
  setFilterTagIds: (tagIds: number[]) => void;
  setSearchQuery: (query: string) => void;
  setPage: (page: number) => void;
  fetchQuestions: () => Promise<void>;
}

export const useQuestionsStore = create<QuestionsState>((set, get) => ({
  items: [],
  totalCount: 0,
  page: 1,
  limit: 20,
  filterDifficulty: '',
  filterType: '',
  filterTagIds: [],
  filterVerified: '',
  searchQuery: '',
  isLoading: false,
  isFetching: false,
  initialized: false,

  setFilter: (key, value) => {
    set({ [key]: value, page: 1 } as Partial<QuestionsState>);
    get().fetchQuestions();
  },

  setFilterTagIds: (tagIds) => {
    set({ filterTagIds: tagIds, page: 1 });
    get().fetchQuestions();
  },

  setSearchQuery: (query) => {
    set({ searchQuery: query, page: 1 });
    get().fetchQuestions();
  },

  setPage: (page) => {
    set({ page });
    get().fetchQuestions();
  },

  fetchQuestions: async () => {
    const { page, limit, filterDifficulty, filterType, filterTagIds, filterVerified, searchQuery, initialized } = get();
    set({ isLoading: !initialized, isFetching: true });
    try {
      const res = await questionsApi.list({
        page,
        limit,
        difficulty: filterDifficulty,
        type: filterType,
        tagIds: filterTagIds.length > 0 ? filterTagIds : undefined,
        verified: filterVerified === '' ? undefined : filterVerified === 'true',
        q: searchQuery || undefined,
      });
      set({ items: res.items ?? [], totalCount: res.total_count, isLoading: false, isFetching: false, initialized: true });
    } catch {
      set({ isLoading: false, isFetching: false });
    }
  },
}));
