import { Injectable, signal } from '@angular/core';
import { ApiService } from './api.service';

type Theme = 'dark' | 'light';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  private readonly storageKey = 'readr_theme';

  theme = signal<Theme>(this.loadTheme());

  constructor(private api: ApiService) {
    this.apply(this.theme());
  }

  toggle(): void {
    const next: Theme = this.theme() === 'dark' ? 'light' : 'dark';
    this.set(next);
  }

  set(theme: Theme): void {
    this.theme.set(theme);
    this.apply(theme);
    localStorage.setItem(this.storageKey, theme);
    this.api.patch('/users/me/settings', { theme }).subscribe({ error: () => {} });
  }

  private apply(theme: Theme): void {
    document.documentElement.classList.toggle('light', theme === 'light');
  }

  private loadTheme(): Theme {
    const stored = localStorage.getItem(this.storageKey);
    return stored === 'light' ? 'light' : 'dark';
  }
}
