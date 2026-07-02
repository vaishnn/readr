import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { Highlight, HighlightColor, HighlightHistory } from '../models';

@Injectable({ providedIn: 'root' })
export class HighlightService {
  constructor(private api: ApiService) {}

  list(bookId: string, page?: number): Observable<Highlight[]> {
    const params = page !== undefined ? { page } : undefined;
    return this.api.get<Highlight[]>(`/books/${bookId}/highlights`, params as Record<string, number>);
  }

  create(bookId: string, payload: { page: number; cfiRange: string; text: string; color: HighlightColor; note?: string }): Observable<Highlight> {
    return this.api.post<Highlight>(`/books/${bookId}/highlights`, payload);
  }

  update(bookId: string, highlightId: string, color: HighlightColor, note: string): Observable<Highlight> {
    return this.api.patch<Highlight>(`/books/${bookId}/highlights/${highlightId}`, { color, note });
  }

  delete(bookId: string, highlightId: string): Observable<void> {
    return this.api.delete<void>(`/books/${bookId}/highlights/${highlightId}`);
  }

  getHistory(bookId: string, highlightId: string): Observable<HighlightHistory[]> {
    return this.api.get<HighlightHistory[]>(`/books/${bookId}/highlights/${highlightId}/history`);
  }
}
