export type BookFormat = 'pdf' | 'epub' | 'cbz';

export interface BookMetadata {
  publisher?: string;
  year?: number;
  language?: string;
  pageCount?: number;
  description?: string;
  isbn?: string;
}

export interface Book {
  id: string;
  ownerId: string;
  title: string;
  author: string;
  format: BookFormat;
  coverUrl?: string;
  metadata: BookMetadata;
  tags: string[];
  isPrivate: boolean;
  allowedUserIds: string[];
  uploadedAt: string;
}

export interface BookListResponse {
  books: Book[];
  total: number;
}

export interface ReadingProgress {
  id: string;
  userId: string;
  bookId: string;
  page: number;
  cfi: string;
  percentage: number;
  zoom: number;
  lastReadAt: string;
  totalSeconds: number;
}

export interface Bookmark {
  id: string;
  userId: string;
  bookId: string;
  page: number;
  cfi: string;
  label: string;
  createdAt: string;
}
