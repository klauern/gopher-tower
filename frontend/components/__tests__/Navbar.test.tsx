import { act, screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { renderWithRouter } from '../../test/utils/test-utils';
import { Navbar } from '../Navbar';

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
    renderWithRouter(<Navbar />);
    expect(screen.getByText('Gopher Tower')).toBeInTheDocument();
  });

  it('renders navigation links', () => {
    renderWithRouter(<Navbar />);
    expect(screen.getByText('Home')).toHaveAttribute('href', '/');
    expect(screen.getByText('Jobs')).toHaveAttribute('href', '/jobs');
  });

  it('applies active styles to current route', () => {
    const { router } = renderWithRouter(<Navbar />);
    router.pathname = '/';

    const homeLink = screen.getByText('Home');
    const jobsLink = screen.getByText('Jobs');

    expect(homeLink.className).toContain('text-blue-600');
    expect(jobsLink.className).toContain('text-gray-500');
  });

  it('applies active styles to jobs route when on jobs page', async () => {
    const { router, rerender } = renderWithRouter(<Navbar />);

    await act(async () => {
      router.pathname = '/jobs';
      rerender(<Navbar />);
    });

    const homeLink = screen.getByRole('link', { name: /home/i });
    const jobsLink = screen.getByRole('link', { name: /jobs/i });

    // Home link should have inactive styles
    expect(homeLink.className).toContain('text-gray-500');
    expect(homeLink.className).toContain('dark:text-gray-300');
    expect(homeLink.className).toContain('hover:text-gray-700');
    expect(homeLink.className).toContain('dark:hover:text-gray-100');

    // Jobs link should have active styles
    expect(jobsLink.className).toContain('text-blue-600');
    expect(jobsLink.className).toContain('dark:text-blue-400');
  });
});
