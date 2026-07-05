// Central registry of all feature flags.
// Add a new entry here to create a new flag — never use raw strings elsewhere.
export const FLAGS = {
  // Allow users to organize books into named collections.
  collections: 'collections',
  // Public unauthenticated browsing of non-private books.
  publicLibrary: 'public-library',
  // Reading time analytics and progress charts.
  readingStats: 'reading-stats',
  // Text highlighting and inline notes.
  highlights: 'highlights',
  // Allow new users to self-register (disable for single-user instances).
  registration: 'registration',
  // Social sharing — shareable links for books and collections.
  socialSharing: 'social-sharing',
  // Global library — browse all public books from any user.
  global: 'global',
  // Popular Books sidebar — shown on every main tab.
  popularBooks: 'popular-books',
} as const;

export type FlagKey = typeof FLAGS[keyof typeof FLAGS];
