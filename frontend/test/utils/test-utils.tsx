import { render } from '@testing-library/react';
import { vi } from 'vitest';

// Mock useRouter
const useRouter = vi.fn();

// Mock next/router
vi.mock('next/router', () => ({
  useRouter: () => useRouter(),
}));

// Default router mock implementation
const createDefaultRouter = () => ({
  route: '/',
  pathname: '/',
  query: {},
  asPath: '/',
  basePath: '',
  push: vi.fn(),
  replace: vi.fn(),
  reload: vi.fn(),
  back: vi.fn(),
  prefetch: vi.fn(),
  beforePopState: vi.fn(),
  events: {
    on: vi.fn(),
    off: vi.fn(),
    emit: vi.fn(),
  },
  isFallback: false,
});

export function renderWithRouter(ui: React.ReactElement) {
  // Set up the router mock for this render
  const router = createDefaultRouter();
  useRouter.mockReturnValue(router);

  const result = render(ui);
  return {
    ...result,
    router,
  };
}
