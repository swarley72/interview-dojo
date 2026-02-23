export interface AuthResponse {
  user_id: string;
  login: string;
  access_token: string;
  refresh_token: string;
}

export interface User {
  id: string;
  login: string;
  is_super_user: boolean;
}

export interface Progress {
  repetitions: number;
  ease_factor: number;
  interval_days: number;
  next_review_at: string;
}

export interface Question {
  id: string;
  title: string;
  content_md: string | null;
  answer_md: string | null;
  excalidraw_json: string | null;
  difficulty: Difficulty;
  type: QuestionType;
  tag_ids: number[];
  verified: boolean;
  progress: Progress | null;
}

export interface QuestionShort {
  id: string;
  title: string;
  difficulty: Difficulty;
  type: QuestionType;
  tag_ids: number[];
  verified: boolean;
}

export interface ListQuestionsResponse {
  total_count: number;
  items: QuestionShort[];
}

export interface Tag {
  id: number;
  name: string;
}

export interface ListTagsResponse {
  tags: Tag[];
}

export interface UserProgressResponse {
  repetitions: number;
  interval_days: number;
  ease_factor: number;
}

export type Difficulty = 'easy' | 'medium' | 'hard';
export type QuestionType = 'theory' | 'coding' | 'algorithm' | 'system_design';
export type AnswerQuality = 'again' | 'hard' | 'good' | 'easy';
