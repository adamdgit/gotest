export type APIError = {
  code: string;
  error: string;
};

export const ApiErrorCodes = {
  UNAUTHORIZED: 'UNAUTHORIZED',
  ACCESS_EXPIRED: 'ACCESS_EXPIRED',
  SESSION_EXPIRED: 'SESSION_EXPIRED',
  INTERNAL_SERVER_ERROR: 'INTERNAL_SERVER_ERROR',
  INVALID_REQUEST: 'INVALID_REQUEST',
  FORBIDDEN: 'FORBIDDEN',
  NOT_FOUND: 'NOT_FOUND',
  CONFLICT: 'CONFLICT',
  RATE_LIMITED: 'RATE_LIMITED',
  NO_ROWS: 'NO_ROWS',
  EMAIL_INUSE: 'EMAIL_INUSE',
  INVALID_BODY: 'INVALID_BODY',
  INVALID_FORM: 'INVALID_FORM',
  INVALID_CREDENTIALS: 'INVALID_CREDENTIALS'
} as const;

export type UserRole =
  | 'admin'
  | 'manager'
  | 'staff'
  | 'member'
  | 'customer';

export interface LoginSessionRes {
  user_id: number;
  access_token: string;
  refresh_token: string;
}

export interface UserSessionRes {
  email: string;
  role: UserRole;
  profile_url: string;
}

export interface AdminUserDataRes {
  id: number;
  email: string;
  firstname: string;
  lastname: string;
  phone: string;
  address: string;
  role: UserRole;
  profile_url: string;
  created_at: string; // ISO string from Go time.Time
  updated_at: string; // ISO string from Go time.Time
}
