import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { Note } from '../models';

@Injectable({ providedIn: 'root' })
export class NoteService {
  constructor(private api: ApiService) {}

  // Pass page to get page-specific notes; omit for book-level notes.
  list(bookId: string, page?: number): Observable<Note[]> {
    const params = page !== undefined ? { page } : undefined;
    return this.api.get<Note[]>(`/books/${bookId}/notes`, params as Record<string, number>);
  }

  create(bookId: string, contentMd: string, page?: number): Observable<Note> {
    return this.api.post<Note>(`/books/${bookId}/notes`, { contentMd, page: page ?? null });
  }

  update(bookId: string, noteId: string, contentMd: string): Observable<Note> {
    return this.api.patch<Note>(`/books/${bookId}/notes/${noteId}`, { contentMd });
  }

  delete(bookId: string, noteId: string): Observable<void> {
    return this.api.delete<void>(`/books/${bookId}/notes/${noteId}`);
  }
}
