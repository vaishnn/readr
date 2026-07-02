import { Component, Input, Output, EventEmitter, OnInit, OnDestroy, ElementRef, ViewChild, AfterViewInit } from '@angular/core';
import { Editor } from '@tiptap/core';
import { StarterKit } from '@tiptap/starter-kit';
import { Highlight } from '@tiptap/extension-highlight';
import { Underline } from '@tiptap/extension-underline';

@Component({
  selector: 'app-tiptap-editor',
  standalone: true,
  templateUrl: './tiptap-editor.component.html',
})
export class TiptapEditorComponent implements AfterViewInit, OnDestroy {
  @Input() initialContent = '';
  @Output() contentChange = new EventEmitter<string>();

  @ViewChild('editorEl') editorEl!: ElementRef<HTMLDivElement>;

  editor!: Editor;

  ngAfterViewInit(): void {
    this.editor = new Editor({
      element: this.editorEl.nativeElement,
      extensions: [StarterKit, Highlight, Underline],
      content: this.initialContent,
      editorProps: {
        attributes: {
          class: 'prose prose-invert prose-sm max-w-none focus:outline-none min-h-[120px] px-3 py-2',
        },
      },
      onUpdate: ({ editor }) => {
        this.contentChange.emit(editor.getHTML());
      },
    });
  }

  ngOnDestroy(): void {
    this.editor?.destroy();
  }

  isActive(mark: string, attrs?: Record<string, unknown>): boolean {
    return attrs ? (this.editor?.isActive(mark, attrs) ?? false) : (this.editor?.isActive(mark) ?? false);
  }

  toggle(command: string): void {
    const chain = this.editor?.chain().focus() as any;
    chain?.[command]?.()?.run();
  }
}
