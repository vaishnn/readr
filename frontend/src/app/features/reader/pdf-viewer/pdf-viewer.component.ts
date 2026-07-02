import { Component, Input, Output, EventEmitter, OnInit } from '@angular/core';
import { NgxExtendedPdfViewerModule } from 'ngx-extended-pdf-viewer';
import { Highlight } from '../../../core/models';

export interface TextSelection {
  text: string;
  page: number;
  cfiRange?: string;
  x: number;
  y: number;
}

const ZOOM_STEP = 1.2;
const MIN_ZOOM = 0.5;
const MAX_ZOOM = 4.0;

@Component({
  selector: 'app-pdf-viewer',
  standalone: true,
  imports: [NgxExtendedPdfViewerModule],
  templateUrl: './pdf-viewer.component.html',
})
export class PdfViewerComponent implements OnInit {
  @Input({ required: true }) src!: string;
  @Input() initialPage = 1;
  @Input() highlights: Highlight[] = [];

  @Output() pageChange         = new EventEmitter<number>();
  @Output() zoomPercentChange  = new EventEmitter<number>();
  @Output() textSelected       = new EventEmitter<TextSelection>();

  authHeaders: Record<string, string> = {};
  viewerZoom: string | number = 'page-width';
  zoomPercent = 100;
  currentPage = 1;

  private currentZoomFactor = 1.0;

  ngOnInit(): void {
    const token = localStorage.getItem('readr_access_token');
    if (token) this.authHeaders = { Authorization: `Bearer ${token}` };
    this.currentPage = this.initialPage || 1;
  }

  zoomIn(): void {
    this.viewerZoom = Math.min(MAX_ZOOM * 100, this.currentZoomFactor * ZOOM_STEP * 100);
  }

  zoomOut(): void {
    this.viewerZoom = Math.max(MIN_ZOOM * 100, this.currentZoomFactor / ZOOM_STEP * 100);
  }

  onPageChange(page: number): void {
    this.currentPage = page;
    this.pageChange.emit(page);
  }

  onCurrentZoomFactor(factor: number): void {
    this.currentZoomFactor = factor;
    this.zoomPercent = Math.round(factor * 100);
    this.zoomPercentChange.emit(this.zoomPercent);
  }

  onMouseUp(): void {
    const sel = window.getSelection();
    if (!sel || sel.isCollapsed) return;
    const text = sel.toString().trim();
    if (!text) return;
    const range = sel.getRangeAt(0);
    const rect = range.getBoundingClientRect();
    this.textSelected.emit({
      text,
      page: this.currentPage,
      x: rect.left + rect.width / 2,
      y: rect.top - 8,
    });
  }
}
