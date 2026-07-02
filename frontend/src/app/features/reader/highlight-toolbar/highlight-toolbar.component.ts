import { Component, Input, Output, EventEmitter } from '@angular/core';
import { NgStyle } from '@angular/common';
import { HighlightColor } from '../../../core/models';

interface ColorOption {
  value: HighlightColor;
  bg: string;
  label: string;
}

@Component({
  selector: 'app-highlight-toolbar',
  standalone: true,
  imports: [NgStyle],
  templateUrl: './highlight-toolbar.component.html',
})
export class HighlightToolbarComponent {
  @Input({ required: true }) position!: { x: number; y: number };
  @Output() colorSelected = new EventEmitter<HighlightColor>();
  @Output() dismiss       = new EventEmitter<void>();

  readonly colors: ColorOption[] = [
    { value: 'yellow', bg: 'bg-yellow-400',  label: 'Yellow' },
    { value: 'green',  bg: 'bg-emerald-400', label: 'Green'  },
    { value: 'blue',   bg: 'bg-blue-400',    label: 'Blue'   },
    { value: 'pink',   bg: 'bg-pink-400',    label: 'Pink'   },
  ];

  get style(): Record<string, string> {
    return {
      position: 'fixed',
      left: `${this.position.x}px`,
      top:  `${this.position.y}px`,
      transform: 'translate(-50%, -100%)',
    };
  }
}
