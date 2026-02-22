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
  isLoading: boolean;

  setFilter: (key: 'filterDifficulty' | 'filterType', value: string) => void;
  setFilterTagIds: (tagIds: number[]) => void;
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
  isLoading: false,

  setFilter: (key, value) => {
    set({ [key]: value, page: 1 } as Partial<QuestionsState>);
    get().fetchQuestions();
  },

  setFilterTagIds: (tagIds) => {
    set({ filterTagIds: tagIds, page: 1 });
    get().fetchQuestions();
  },

  setPage: (page) => {
    set({ page });
    get().fetchQuestions();
  },

  fetchQuestions: async () => {
    const { page, limit, filterDifficulty, filterType, filterTagIds } = get();
    set({ isLoading: true });
    try {
      const res = await questionsApi.list({
        page,
        limit,
        difficulty: filterDifficulty,
        type: filterType,
        tagIds: filterTagIds.length > 0 ? filterTagIds : undefined,
      });
      set({ items: res.items ?? [], totalCount: res.total_count, isLoading: false });
    } catch {
      set({ isLoading: false });
    }
  },
}));
