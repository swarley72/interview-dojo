import { create } from 'zustand';
import type { Tag } from '../types';
import { tagsApi } from '../api/tags';

interface TagsState {
  tags: Tag[];
  tagMap: Record<number, string>;
  fetchTags: () => Promise<void>;
}

export const useTagsStore = create<TagsState>((set) => ({
  tags: [],
  tagMap: {},

  fetchTags: async () => {
    try {
      const res = await tagsApi.list();
      const tags = res.tags ?? [];
      const tagMap: Record<number, string> = {};
      for (const t of tags) tagMap[t.id] = t.name;
      set({ tags, tagMap });
    } catch {
      // ignore
    }
  },
}));
