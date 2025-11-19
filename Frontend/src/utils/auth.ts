// Utility functions for authentication

export interface CurrentUser {
  user_uid: string;
  email: string;
  name: string;
}

export function getCurrentUser(): CurrentUser | null {
  const userStr = localStorage.getItem('user');
  if (!userStr) return null;
  
  try {
    return JSON.parse(userStr);
  } catch (error) {
    console.error('Failed to parse user from localStorage:', error);
    return null;
  }
}

export function getCurrentUserUid(): string | null {
  const user = getCurrentUser();
  return user?.user_uid || null;
}

export function isUserLoggedIn(): boolean {
  return !!localStorage.getItem('token') && !!getCurrentUser();
}
