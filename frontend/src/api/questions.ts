import type { ListQuestionsResponse, Question, Difficulty, QuestionType } from '../types';
import { api } from './client';

interface ListParams {
  page?: number;
  limit?: number;
  difficulty?: Difficulty | '';
  type?: QuestionType | '';
  tagIds?: number[];
}

export const questionsApi = {
  list: (params: ListParams = {}) => {
    const query = new URLSearchParams();
    if (params.page) query.set('page', String(params.page));
    if (params.limit) query.set('limit', String(params.limit));
    if (params.difficulty) query.set('difficulty', params.difficulty);
    if (params.type) query.set('type', params.type);
    if (params.tagIds) {
      for (const id of params.tagIds) query.append('tag_id', String(id));
    }
    return api.get<ListQuestionsResponse>(`/questions?${query}`);
  },

  get: (id: string) => api.get<Question>(`/questions/${id}`),

  create: (data: {
    title: string;
    content_md?: string;
    answer_md?: string;
    excalidraw_json?: string;
    difficulty: Difficulty;
    type: QuestionType;
    tag_ids?: number[];
  }) => api.post<Question>('/questions', data),

  update: (
    id: string,
    data: {
      title?: string;
      content_md?: string;
      answer_md?: string;
      excalidraw_json?: string;
      difficulty?: Difficulty;
      type?: QuestionType;
      tag_ids?: number[];
    },
  ) => api.patch<Question>(`/questions/${id}`, data),

  delete: (id: string) => api.del(`/questions/${id}`),
};
