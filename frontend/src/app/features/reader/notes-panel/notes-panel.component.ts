import { Component, Input, Output, EventEmitter, OnInit, signal } from '@angular/core';
import { Book, Highlight, Note } from '../../../core/models';
import { NoteService } from '../../../core/services/note.service';
import { HighlightService } from '../../../core/services/highlight.service';
import { ToastService } from '../../../shared/components/toast.service';
import { TiptapEditorComponent } from './tiptap-editor.component';
import { SpinnerComponent } from '../../../shared/components/spinner.component';

type Tab = 'book-notes' | 'page-notes' | 'highlights';

@Component({
  selector: 'app-notes-panel',
  standalone: true,
  imports: [TiptapEditorComponent, SpinnerComponent],
  templateUrl: './notes-panel.component.html',
})
export class NotesPanelComponent implements OnInit {
  @Input({ required: true }) book!: Book;
  @Input({ required: true }) currentPage!: number;
  @Input() highlights: Highlight[] = [];

  @Output() highlightUpdated = new EventEmitter<Highlight>();
  @Output() highlightDeleted = new EventEmitter<string>();
  @Output() close            = new EventEmitter<void>();

  activeTab  = signal<Tab>('book-notes');
  bookNotes  = signal<Note[]>([]);
  pageNotes  = signal<Note[]>([]);
  loading    = signal(true);
  savingNote = signal(false);

  readonly tabs: { id: Tab; label: string }[] = [
    { id: 'book-notes',  label: 'Book'       },
    { id: 'page-notes',  label: 'This page'  },
    { id: 'highlights',  label: 'Highlights' },
  ];

  readonly highlightColors: { [key: string]: string } = {
    yellow: 'bg-yellow-400/30 border-yellow-400',
    green:  'bg-emerald-400/30 border-emerald-400',
    blue:   'bg-blue-400/30 border-blue-400',
    pink:   'bg-pink-400/30 border-pink-400',
  };

  constructor(
    private noteService: NoteService,
    private highlightService: HighlightService,
    private toast: ToastService,
  ) {}

  ngOnInit(): void {
    this.loadNotes();
  }

  loadNotes(): void {
    this.loading.set(true);
    this.noteService.list(this.book.id).subscribe({
      next: notes => {
        this.bookNotes.set(notes);
        this.loading.set(false);
      },
    });
    this.noteService.list(this.book.id, this.currentPage).subscribe({
      next: notes => this.pageNotes.set(notes),
    });
  }

  saveBookNote(content: string): void {
    const existing = this.bookNotes()[0];
    const req$ = existing
      ? this.noteService.update(this.book.id, existing.id, content)
      : this.noteService.create(this.book.id, content);

    req$.subscribe({
      next: note => {
        this.bookNotes.set([note]);
      },
      error: () => this.toast.error('Failed to save note'),
    });
  }

  savePageNote(content: string): void {
    const existing = this.pageNotes()[0];
    const req$ = existing
      ? this.noteService.update(this.book.id, existing.id, content)
      : this.noteService.create(this.book.id, content, this.currentPage);

    req$.subscribe({
      next: note => {
        this.pageNotes.set([note]);
      },
      error: () => this.toast.error('Failed to save note'),
    });
  }

  bookNoteContent(): string {
    const notes = this.bookNotes();
    return notes.length > 0 ? notes[0].contentMd : '';
  }

  pageNoteContent(): string {
    const notes = this.pageNotes();
    return notes.length > 0 ? notes[0].contentMd : '';
  }

  highlightColorClass(color: string): string {
    return this.highlightColors[color] ?? 'bg-slate-700/30 border-slate-600';
  }

  deleteHighlight(id: string): void {
    this.highlightService.delete(this.book.id, id).subscribe({
      next: () => {
        this.highlightDeleted.emit(id);
        this.toast.success('Highlight removed');
      },
      error: () => this.toast.error('Failed to delete highlight'),
    });
  }
}
