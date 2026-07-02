import { Component, Input, Output, EventEmitter, OnInit, OnDestroy, AfterViewInit, ViewChild, ElementRef } from '@angular/core';
import { Highlight } from '../../../core/models';
import { TextSelection } from '../pdf-viewer/pdf-viewer.component';

@Component({
  selector: 'app-epub-viewer',
  standalone: true,
  templateUrl: './epub-viewer.component.html',
})
export class EpubViewerComponent implements AfterViewInit, OnDestroy {
  @Input({ required: true }) src!: string;
  @Input() initialCFI = '';
  @Input() highlights: Highlight[] = [];

  @Output() cfiChange    = new EventEmitter<string>();
  @Output() textSelected = new EventEmitter<TextSelection>();

  @ViewChild('viewer') viewerRef!: ElementRef<HTMLDivElement>;

  loading = true;

  private book: any     = null;
  private rendition: any = null;

  async ngAfterViewInit(): Promise<void> {
    const ePub = (await import('epubjs')).default;

    this.book = ePub(this.src);
    this.rendition = this.book.renderTo(this.viewerRef.nativeElement, {
      width:  '100%',
      height: '100%',
      spread: 'none',
    });

    if (this.initialCFI) {
      await this.rendition.display(this.initialCFI);
    } else {
      await this.rendition.display();
    }

    this.loading = false;

    this.rendition.on('relocated', (location: any) => {
      this.cfiChange.emit(location.start.cfi);
    });

    // Wire up text selection inside the epub iframe.
    this.rendition.on('selected', (cfiRange: string, contents: any) => {
      const sel  = contents.window.getSelection();
      const text = sel?.toString().trim();
      if (!text) return;

      const range = sel?.getRangeAt(0);
      const rect  = range?.getBoundingClientRect();
      if (!rect) return;

      // Offset by the iframe position within the page.
      const iframeRect = this.viewerRef.nativeElement.querySelector('iframe')?.getBoundingClientRect();
      const x = (iframeRect?.left ?? 0) + rect.left + rect.width / 2;
      const y = (iframeRect?.top  ?? 0) + rect.top - 8;

      this.textSelected.emit({ text, page: 0, cfiRange, x, y });
    });
  }

  ngOnDestroy(): void {
    this.book?.destroy();
  }

  prevPage(): void  { this.rendition?.prev(); }
  nextPage(): void  { this.rendition?.next(); }
}
