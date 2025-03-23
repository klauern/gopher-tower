import { AppRouterContext, AppRouterInstance } from 'next/dist/shared/lib/app-router-context.shared-runtime';
import { ReactNode } from 'react';

interface MockRouterContextProps {
  children: ReactNode;
  push?: (href: string) => void;
}

export function MockRouterContext({ children, push = () => {} }: MockRouterContextProps) {
  const mockRouter: AppRouterInstance = {
    back: () => {},
    forward: () => {},
    push,
    refresh: () => {},
    replace: () => {},
    prefetch: () => Promise.resolve(),
  };

  return (
    <AppRouterContext.Provider value={mockRouter}>
      {children}
    </AppRouterContext.Provider>
  );
}

export function renderWithRouter(ui: ReactNode, { push }: { push?: (href: string) => void } = {}) {
  return (
    <MockRouterContext push={push}>
      {ui}
    </MockRouterContext>
  );
}
