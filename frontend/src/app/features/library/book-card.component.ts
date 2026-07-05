import { Component, Input, Output, EventEmitter, signal, computed, HostListener } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { DatePipe } from '@angular/common';
import { Book } from '../../core/models';
import { BookService } from '../../core/services/book.service';
import { AuthService } from '../../core/services/auth.service';
import { ToastService } from '../../shared/components/toast.service';

@Component({
  selector: 'app-book-card',
  standalone: true,
  imports: [FormsModule, DatePipe],
  templateUrl: './book-card.component.html',
})
export class BookCardComponent {
  @Input({ required: true }) book!: Book;
  @Output() open    = new EventEmitter<Book>();
  @Output() deleted = new EventEmitter<string>();
  @Output() updated = new EventEmitter<Book>();

  menuOpen = signal(false);
  deleting = signal(false);
  editOpen = signal(false);
  metaOpen = signal(false);
  saving   = signal(false);

  // Edit form fields
  editTitle       = '';
  editAuthor      = '';
  editTags        = '';
  editPublisher   = '';
  editYear: number | null = null;
  editLanguage    = '';
  editIsbn        = '';
  editDescription = '';

  isOwner = computed(() => this.auth.currentUser()?.id === this.book.ownerId);

  constructor(
    private bookService: BookService,
    private auth: AuthService,
    private toast: ToastService,
  ) {}

  get formatBadgeClass(): string {
    return {
      pdf:  'bg-red-900/60 text-red-300',
      epub: 'bg-blue-900/60 text-blue-300',
      cbz:  'bg-purple-900/60 text-purple-300',
    }[this.book.format] ?? 'bg-slate-700 text-slate-300';
  }

  toggleMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.menuOpen.update(v => !v);
  }

  openEdit(event: MouseEvent): void {
    event.stopPropagation();
    this.menuOpen.set(false);
    this.editTitle       = this.book.title;
    this.editAuthor      = this.book.author ?? '';
    this.editTags        = (this.book.tags ?? []).join(', ');
    this.editPublisher   = this.book.metadata.publisher ?? '';
    this.editYear        = this.book.metadata.year ?? null;
    this.editLanguage    = this.book.metadata.language ?? '';
    this.editIsbn        = this.book.metadata.isbn ?? '';
    this.editDescription = this.book.metadata.description ?? '';
    this.editOpen.set(true);
  }

  closeEdit(): void {
    this.editOpen.set(false);
  }

  saveEdit(): void {
    if (!this.editTitle.trim()) return;
    this.saving.set(true);
    const tags = this.editTags.split(',').map(t => t.trim()).filter(Boolean);
    this.bookService.update(this.book.id, {
      title: this.editTitle.trim(),
      author: this.editAuthor.trim(),
      tags,
      metadata: {
        publisher: this.editPublisher || undefined,
        year: this.editYear || undefined,
        language: this.editLanguage || undefined,
        isbn: this.editIsbn || undefined,
        description: this.editDescription || undefined,
        pageCount: this.book.metadata.pageCount,
      },
    }).subscribe({
      next: saved => {
        this.saving.set(false);
        this.editOpen.set(false);
        this.toast.success('Book updated');
        this.updated.emit(saved);
      },
      error: err => {
        this.saving.set(false);
        this.toast.error(err?.error?.error ?? 'Failed to update book');
      },
    });
  }

  openMeta(event: MouseEvent): void {
    event.stopPropagation();
    this.menuOpen.set(false);
    this.metaOpen.set(true);
  }

  delete(event: MouseEvent): void {
    event.stopPropagation();
    if (!confirm(`Delete "${this.book.title}"?`)) return;
    this.deleting.set(true);
    this.bookService.delete(this.book.id).subscribe({
      next: () => {
        this.toast.success('Book deleted');
        this.deleted.emit(this.book.id);
      },
      error: err => {
        this.deleting.set(false);
        this.toast.error(err?.error?.error ?? 'Failed to delete book');
      },
    });
  }

  @HostListener('document:click')
  closeMenu(): void {
    this.menuOpen.set(false);
  }
}
