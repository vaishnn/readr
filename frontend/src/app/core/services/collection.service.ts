import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ApiService } from './api.service';
import { Collection } from '../models';

@Injectable({ providedIn: 'root' })
export class CollectionService {
  constructor(private api: ApiService) {}

  list(): Observable<Collection[]> {
    return this.api.get<Collection[]>('/collections');
  }

  create(name: string, description: string): Observable<Collection> {
    return this.api.post<Collection>('/collections', { name, description });
  }

  update(id: string, name: string, description: string): Observable<Collection> {
    return this.api.patch<Collection>(`/collections/${id}`, { name, description });
  }

  delete(id: string): Observable<void> {
    return this.api.delete<void>(`/collections/${id}`);
  }

  addBook(collectionId: string, bookId: string): Observable<void> {
    return this.api.post<void>(`/collections/${collectionId}/books`, { bookId });
  }

  removeBook(collectionId: string, bookId: string): Observable<void> {
    return this.api.delete<void>(`/collections/${collectionId}/books/${bookId}`);
  }
}
