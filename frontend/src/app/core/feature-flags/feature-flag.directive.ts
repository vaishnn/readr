import { Directive, Input, OnInit, TemplateRef, ViewContainerRef } from '@angular/core';
import { FeatureFlagService } from './feature-flag.service';
import { FlagKey } from './flags';

/**
 * Structural directive — renders the host element only when the flag is enabled.
 *
 * Usage:
 *   <div *appFeatureFlag="'collections'">...</div>
 *   <app-collections *appFeatureFlag="'collections'" />
 */
@Directive({
  selector: '[appFeatureFlag]',
  standalone: true,
})
export class FeatureFlagDirective implements OnInit {
  @Input({ required: true }) appFeatureFlag!: FlagKey;

  constructor(
    private template: TemplateRef<unknown>,
    private vcr: ViewContainerRef,
    private flags: FeatureFlagService,
  ) {}

  ngOnInit(): void {
    if (this.flags.isEnabled(this.appFeatureFlag)) {
      this.vcr.createEmbeddedView(this.template);
    }
  }
}
