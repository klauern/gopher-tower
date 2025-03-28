import { screen } from '@testing-library/react';
import { describe, expect, it, vi } from 'vitest';
import { renderWithRouter } from '../../test/utils/test-utils';
import { Navbar } from '../Navbar';

const usePathname = vi.fn();

// Mock next/navigation
vi.mock('next/navigation', () => ({
  usePathname: () => usePathname()
}));

describe('Navbar', () => {
  it('renders the logo text', () => {
    usePathname.mockReturnValue('/');
    renderWithRouter(<Navbar />);
    expect(screen.getByText('Gopher Tower')).toBeInTheDocument();
  });

  it('renders navigation links', () => {
    usePathname.mockReturnValue('/');
    renderWithRouter(<Navbar />);

    const homeLink = screen.getByRole('link', { name: /home/i });
    const jobsLink = screen.getByRole('link', { name: /jobs/i });
    const eventsLink = screen.getByRole('link', { name: /events/i });

    expect(homeLink).toBeInTheDocument();
    expect(jobsLink).toBeInTheDocument();
    expect(eventsLink).toBeInTheDocument();
  });

  it('applies active styles to current route', () => {
    usePathname.mockReturnValue('/');
    renderWithRouter(<Navbar />);

    const homeLink = screen.getByRole('link', { name: /home/i });
    const jobsLink = screen.getByRole('link', { name: /jobs/i });
    const eventsLink = screen.getByRole('link', { name: /events/i });

    expect(homeLink.className).toContain('!text-blue-600');
    expect(jobsLink.className).toContain('text-gray-500');
    expect(eventsLink.className).toContain('text-gray-500');
  });

  it('applies active styles to jobs route when on jobs page', () => {
    usePathname.mockReturnValue('/jobs');
    renderWithRouter(<Navbar />);

    const homeLink = screen.getByRole('link', { name: /home/i });
    const jobsLink = screen.getByRole('link', { name: /jobs/i });
    const eventsLink = screen.getByRole('link', { name: /events/i });

    // Home link should have inactive styles
    expect(homeLink.className).toContain('text-gray-500');
    expect(homeLink.className).toContain('dark:text-gray-300');
    expect(homeLink.className).toContain('hover:text-gray-700');

    // Jobs link should have active styles
    expect(jobsLink.className).toContain('!text-blue-600');
    expect(jobsLink.className).toContain('dark:!text-blue-400');

    // Events link should have inactive styles
    expect(eventsLink.className).toContain('text-gray-500');
    expect(eventsLink.className).toContain('dark:text-gray-300');
    expect(eventsLink.className).toContain('hover:text-gray-700');
  });

  it('applies active styles to events route when on events page', () => {
    usePathname.mockReturnValue('/events');
    renderWithRouter(<Navbar />);

    const homeLink = screen.getByRole('link', { name: /home/i });
    const jobsLink = screen.getByRole('link', { name: /jobs/i });
    const eventsLink = screen.getByRole('link', { name: /events/i });

    // Home and Jobs links should have inactive styles
    expect(homeLink.className).toContain('text-gray-500');
    expect(jobsLink.className).toContain('text-gray-500');

    // Events link should have active styles
    expect(eventsLink.className).toContain('!text-blue-600');
    expect(eventsLink.className).toContain('dark:!text-blue-400');
  });
});
