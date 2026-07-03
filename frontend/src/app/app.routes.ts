import { Routes } from '@angular/router';
import { authGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  { path: '', redirectTo: 'library', pathMatch: 'full' },
  {
    path: 'login',
    loadComponent: () =>
      import('./features/auth/auth.component').then(m => m.AuthComponent),
  },
  {
    path: 'register',
    loadComponent: () =>
      import('./features/auth/auth.component').then(m => m.AuthComponent),
  },
  // Full-screen reader — no shell wrapper
  {
    path: 'reader/:bookId',
    canActivate: [authGuard],
    loadComponent: () =>
      import('./features/reader/reader.component').then(m => m.ReaderComponent),
  },
  // Shell-wrapped authenticated routes
  {
    path: '',
    loadComponent: () =>
      import('./shared/components/shell.component').then(m => m.ShellComponent),
    canActivate: [authGuard],
    children: [
      {
        path: 'library',
        loadComponent: () =>
          import('./features/library/library.component').then(m => m.LibraryComponent),
      },
      {
        path: 'collections',
        loadComponent: () =>
          import('./features/collections/collections.component').then(m => m.CollectionsComponent),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./features/settings/settings.component').then(m => m.SettingsComponent),
      },
    ],
  },
  { path: '**', redirectTo: 'library' },
];
