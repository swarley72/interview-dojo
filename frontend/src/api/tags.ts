import type { ListTagsResponse, Tag } from '../types';
import { api } from './client';

export const tagsApi = {
  list: () => api.get<ListTagsResponse>('/tags'),
  create: (name: string) => api.post<Tag>('/tags', { name }),
  delete: (id: number) => api.del(`/tags/${id}`),
};
