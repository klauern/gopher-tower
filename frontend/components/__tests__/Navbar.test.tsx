import { render, screen } from '@testing-library/react';
import { useRouter } from 'next/router';
import { describe, expect, it, vi } from 'vitest';
import { Navbar } from '../Navbar';

// Mock next/router
vi.mock('next/router', () => ({
  useRouter: () => ({
    pathname: '/',
    route: '/',
    query: {},
    asPath: '/',
    basePath: '',
    isLocaleDomain: false,
    push: vi.fn(),
    replace: vi.fn(),
    reload: vi.fn(),
    back: vi.fn(),
    forward: vi.fn(),
    prefetch: vi.fn(),
    beforePopState: vi.fn(),
    events: {
      on: vi.fn(),
      off: vi.fn(),
      emit: vi.fn(),
    },
    isFallback: false,
    isReady: true,
    isPreview: false,
  }),
}));

// Mock next/link
vi.mock('next/link', () => {
  return {
    __esModule: true,
    default: ({ children, href, className }: { children: React.ReactNode; href: string; className: string }) => (
      <a href={href} className={className}>
        {children}
      </a>
    ),
  };
});

describe('Navbar', () => {
  it('renders the logo text', () => {
    render(<Navbar />);
    expect(screen.getByText('Gopher Tower')).toBeInTheDocument();
  });

  it('renders navigation links', () => {
    render(<Navbar />);
    expect(screen.getByText('Home')).toHaveAttribute('href', '/');
    expect(screen.getByText('Jobs')).toHaveAttribute('href', '/jobs');
  });

  it('applies active styles to current route', () => {
    render(<Navbar />);
    const homeLink = screen.getByText('Home');
    const jobsLink = screen.getByText('Jobs');

    // Home should be active (mocked pathname is '/')
    expect(homeLink.className).toContain('text-blue-600');
    expect(jobsLink.className).toContain('text-gray-500');
  });

  it('applies active styles to jobs route when on jobs page', () => {
    // Update mock to simulate being on jobs page
    vi.mocked(useRouter).mockReturnValue({
      pathname: '/jobs',
      route: '/jobs',
      query: {},
      asPath: '/jobs',
      basePath: '',
      isLocaleDomain: false,
      push: vi.fn(),
      replace: vi.fn(),
      reload: vi.fn(),
      back: vi.fn(),
      forward: vi.fn(),
      prefetch: vi.fn(),
      beforePopState: vi.fn(),
      events: {
        on: vi.fn(),
        off: vi.fn(),
        emit: vi.fn(),
      },
      isFallback: false,
      isReady: true,
      isPreview: false,
    });

    render(<Navbar />);
    const homeLink = screen.getByText('Home');
    const jobsLink = screen.getByText('Jobs');

    expect(homeLink.className).toContain('text-gray-500');
    expect(jobsLink.className).toContain('text-blue-600');
  });
});
