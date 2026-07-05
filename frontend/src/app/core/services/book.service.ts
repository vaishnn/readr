import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { Book, BookListResponse, BookMetadata, ReadingProgress, Bookmark } from '../models';

export interface BookListParams {
  page?: number;
  limit?: number;
  search?: string;
  tag?: string;
  ownerOnly?: boolean;
}

@Injectable({ providedIn: 'root' })
export class BookService {
  constructor(private api: ApiService) {}

  list(params?: BookListParams): Observable<BookListResponse> {
    return this.api.get<BookListResponse>('/books', params as Record<string, string | number>);
  }

  get(id: string): Observable<Book> {
    return this.api.get<Book>(`/books/${id}`);
  }

  upload(file: File, cover?: File, fields?: { title?: string; author?: string; tags?: string[]; metadata?: Partial<BookMetadata> }): Observable<Book> {
    const form = new FormData();
    form.append('book', file);
    if (cover)                form.append('cover',    cover);
    if (fields?.title)        form.append('title',    fields.title);
    if (fields?.author)       form.append('author',   fields.author);
    if (fields?.tags?.length) form.append('tags',     JSON.stringify(fields.tags));
    if (fields?.metadata)     form.append('metadata', JSON.stringify(fields.metadata));
    return this.api.upload<Book>('/books', form);
  }

  update(id: string, data: { title: string; author: string; tags: string[]; metadata: Partial<BookMetadata> }): Observable<Book> {
    return this.api.patch<Book>(`/books/${id}`, data);
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/books/${id}`);
  }

  // Returns a presigned redirect URL — use directly as an <iframe> or fetch src.
  streamUrl(id: string): string {
    return `/api/v1/books/${id}/stream`;
  }

  updateAccess(id: string, isPrivate: boolean, allowedUserIds: string[]): Observable<void> {
    return this.api.patch<void>(`/books/${id}/access`, { isPrivate, allowedUserIds });
  }

  getProgress(bookId: string): Observable<ReadingProgress | null> {
    return this.api.get<ReadingProgress | null>(`/books/${bookId}/progress`);
  }

  saveProgress(bookId: string, page: number, cfi: string, percentage: number, zoom: number, sessionSeconds: number): Observable<void> {
    return this.api.put<void>(`/books/${bookId}/progress`, { page, cfi, percentage, zoom, sessionSeconds });
  }

  listBookmarks(bookId: string): Observable<Bookmark[]> {
    return this.api.get<Bookmark[]>(`/books/${bookId}/bookmarks`);
  }

  createBookmark(bookId: string, page: number, cfi: string, label: string): Observable<Bookmark> {
    return this.api.post<Bookmark>(`/books/${bookId}/bookmarks`, { page, cfi, label });
  }

  deleteBookmark(bookId: string, bookmarkId: string): Observable<void> {
    return this.api.delete<void>(`/books/${bookId}/bookmarks/${bookmarkId}`);
  }
}
