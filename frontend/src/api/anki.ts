import type { AnswerQuality, Question, UserProgressResponse } from '../types';
import { api } from './client';

export const ankiApi = {
  getNext: (params?: { type?: string; tagIds?: number[] }) => {
    const q = new URLSearchParams();
    if (params?.type) q.set('type', params.type);
    params?.tagIds?.forEach((id) => q.append('tag_id', String(id)));
    const qs = q.toString();
    return api.get<Question>(`/anki/next${qs ? `?${qs}` : ''}`);
  },
  recordAnswer: (questionId: string, answer: AnswerQuality) =>
    api.post<UserProgressResponse>(`/anki/${questionId}/answer`, { answer }),
};
