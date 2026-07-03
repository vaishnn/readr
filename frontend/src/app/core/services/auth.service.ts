import { Injectable, signal, computed } from '@angular/core';
import { Router } from '@angular/router';
import { Observable, tap } from 'rxjs';
import { ApiService } from './api.service';
import { AuthResponse, User, TokenPair } from '../models';

const ACCESS_TOKEN_KEY = 'readr_access_token';
const REFRESH_TOKEN_KEY = 'readr_refresh_token';

@Injectable({ providedIn: 'root' })
export class AuthService {
  private _currentUser = signal<User | null>(this.loadUser());

  readonly currentUser = this._currentUser.asReadonly();
  readonly isLoggedIn = computed(() => this._currentUser() !== null);

  constructor(private api: ApiService, private router: Router) {}

  login(email: string, password: string): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/auth/login', { email, password }).pipe(
      tap(res => this.persistSession(res))
    );
  }

  register(email: string, username: string, password: string): Observable<AuthResponse> {
    return this.api.post<AuthResponse>('/auth/register', { email, username, password }).pipe(
      tap(res => this.persistSession(res))
    );
  }

  refresh(): Observable<TokenPair> {
    const refreshToken = localStorage.getItem(REFRESH_TOKEN_KEY);
    return this.api.post<TokenPair>('/auth/refresh', { refreshToken }).pipe(
      tap(tokens => {
        localStorage.setItem(ACCESS_TOKEN_KEY, tokens.accessToken);
        localStorage.setItem(REFRESH_TOKEN_KEY, tokens.refreshToken);
      })
    );
  }

  logout(): void {
    this.api.delete('/auth/logout').subscribe({ error: () => {} });
    this.clearSession();
    this.router.navigate(['/login']);
  }

  getAccessToken(): string | null {
    return localStorage.getItem(ACCESS_TOKEN_KEY);
  }

  getRefreshToken(): string | null {
    return localStorage.getItem(REFRESH_TOKEN_KEY);
  }

  updateUser(user: User): void {
    localStorage.setItem('readr_user', JSON.stringify(user));
    this._currentUser.set(user);
  }

  private persistSession(res: AuthResponse): void {
    localStorage.setItem(ACCESS_TOKEN_KEY, res.tokens.accessToken);
    localStorage.setItem(REFRESH_TOKEN_KEY, res.tokens.refreshToken);
    localStorage.setItem('readr_user', JSON.stringify(res.user));
    this._currentUser.set(res.user);
  }

  private clearSession(): void {
    localStorage.removeItem(ACCESS_TOKEN_KEY);
    localStorage.removeItem(REFRESH_TOKEN_KEY);
    localStorage.removeItem('readr_user');
    this._currentUser.set(null);
  }

  private loadUser(): User | null {
    try {
      const raw = localStorage.getItem('readr_user');
      return raw ? JSON.parse(raw) : null;
    } catch {
      return null;
    }
  }
}
