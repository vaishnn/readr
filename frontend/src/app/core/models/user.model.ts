export interface UserSettings {
  theme: 'dark' | 'light';
  defaultView: 'grid' | 'list';
  librarySidebarOpen: boolean;
}

export interface User {
  id: string;
  email: string;
  username: string;
  settings: UserSettings;
  createdAt: string;
  updatedAt: string;
}

export interface TokenPair {
  accessToken: string;
  refreshToken: string;
}

export interface AuthResponse {
  user: User;
  tokens: TokenPair;
}
