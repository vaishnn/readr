import { Injectable } from '@angular/core';

@Injectable({ providedIn: 'root' })
export class ThemeService {
  constructor() {
    document.documentElement.classList.remove('light');
    localStorage.removeItem('readr_theme');
  }
}
