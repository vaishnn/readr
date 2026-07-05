import { Component, OnInit, signal } from '@angular/core';
import { Router } from '@angular/router';
import { BookService } from '../../core/services/book.service';
import { AuthService } from '../../core/services/auth.service';
import { ApiService } from '../../core/services/api.service';
import { Book, User } from '../../core/models';

@Component({
  selector: 'app-popular-books',
  standalone: true,
  imports: [],
  templateUrl: './popular-books.component.html',
})
export class PopularBooksComponent implements OnInit {
  books       = signal<Book[]>([]);
  sidebarOpen = signal(true);

  constructor(
    private bookService: BookService,
    private router: Router,
    private auth: AuthService,
    private api: ApiService,
  ) {}

  ngOnInit(): void {
    this.sidebarOpen.set(this.auth.currentUser()?.settings?.librarySidebarOpen ?? true);
    this.bookService.list({ page: 1, limit: 6, ownerOnly: true }).subscribe({
      next: res => this.books.set(res.books ?? []),
      error: () => {},
    });
  }

  toggleSidebar(): void {
    const next = !this.sidebarOpen();
    this.sidebarOpen.set(next);
    const user = this.auth.currentUser();
    if (!user) return;
    const updated: User = { ...user, settings: { ...user.settings, librarySidebarOpen: next } };
    this.auth.updateUser(updated);
    this.api.patch<User>('/users/me/settings', updated.settings).subscribe({
      next: saved => this.auth.updateUser(saved),
      error: () => {},
    });
  }

  openBook(book: Book): void {
    this.router.navigate(['/reader', book.id]);
  }
}
