import { Injectable, signal } from '@angular/core';
import { firstValueFrom, interval } from 'rxjs';
import { switchMap } from 'rxjs/operators';
import { ApiService } from '../services/api.service';
import { FlagKey } from './flags';

type FlagMap = Record<string, boolean>;

const POLL_INTERVAL_MS = 10_000;

@Injectable({ providedIn: 'root' })
export class FeatureFlagService {
  private serverFlags = signal<FlagMap>({});

  constructor(private api: ApiService) {}

  // Called by APP_INITIALIZER — Angular blocks all rendering until this resolves.
  // On failure the app still boots (flags default to false) rather than hard-crashing.
  async init(): Promise<void> {
    try {
      const flags = await firstValueFrom(this.api.get<FlagMap>('/features'));
      this.serverFlags.set(flags);
    } catch {
      // Server unreachable — app boots with all flags false.
    }
    this.startPolling();
  }

  isEnabled(key: FlagKey): boolean {
    return this.serverFlags()[key] ?? false;
  }

  isEnabled$(key: FlagKey): () => boolean {
    return () => this.isEnabled(key);
  }

  private startPolling(): void {
    interval(POLL_INTERVAL_MS)
      .pipe(switchMap(() => this.api.get<FlagMap>('/features')))
      .subscribe({
        next: flags => {
          if (this.flagsChanged(flags)) {
            window.location.reload();
          }
        },
        error: () => {},
      });
  }

  private flagsChanged(incoming: FlagMap): boolean {
    const current = this.serverFlags();
    const keys = new Set([...Object.keys(current), ...Object.keys(incoming)]);
    for (const key of keys) {
      if (current[key] !== incoming[key]) return true;
    }
    return false;
  }
}
