import { describe, expect, it, vi } from 'vitest';
import Page from '../page';
import { render, screen } from '../utils/test-utils';

// Mock the EventStream component since we're not testing its functionality here
vi.mock('../components/EventStream', () => ({
  EventStream: ({ children }: { children: (isConnected: boolean) => React.ReactNode }) =>
    children ? children(false) : null
}));

describe('Page Component', () => {
  it('renders the page with expected elements', () => {
    render(<Page />);

    // Check for key elements that should be present on the page
    // Adjust these assertions based on what's actually in your page component
    expect(screen.getByRole('heading', { level: 1 })).toBeDefined();
  });
});
