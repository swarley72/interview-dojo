import type { AnswerQuality, Question, UserProgressResponse } from '../types';
import { api } from './client';

export const ankiApi = {
  getNext: () => api.get<Question>('/anki/next'),
  recordAnswer: (questionId: string, answer: AnswerQuality) =>
    api.post<UserProgressResponse>(`/anki/${questionId}/answer`, { answer }),
};
