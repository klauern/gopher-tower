import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { describe, expect, it, vi } from 'vitest';
import { Button } from '../Button';

describe('Button', () => {
  it('renders with default props', () => {
    render(<Button>Click me</Button>);
    const button = screen.getByText('Click me');

    expect(button).toBeInTheDocument();
    expect(button).toHaveClass('bg-primary');
    expect(button).toHaveClass('h-9 px-4 py-2'); // default size
  });

  it('renders with different variants', () => {
    const { rerender } = render(<Button variant="secondary">Secondary</Button>);
    let button = screen.getByText('Secondary');
    expect(button).toHaveClass('bg-secondary');

    rerender(<Button variant="destructive">Destructive</Button>);
    button = screen.getByText('Destructive');
    expect(button).toHaveClass('bg-destructive');

    rerender(<Button variant="outline">Outline</Button>);
    button = screen.getByText('Outline');
    expect(button).toHaveClass('border-input');
  });

  it('renders with different sizes', () => {
    const { rerender } = render(<Button size="sm">Small</Button>);
    let button = screen.getByText('Small');
    expect(button).toHaveClass('h-8');
    expect(button).toHaveClass('px-3');
    expect(button).toHaveClass('text-xs');

    rerender(<Button size="lg">Large</Button>);
    button = screen.getByText('Large');
    expect(button).toHaveClass('h-10');
    expect(button).toHaveClass('px-8');

    rerender(<Button size="icon">Icon</Button>);
    button = screen.getByText('Icon');
    expect(button).toHaveClass('h-9');
    expect(button).toHaveClass('w-9');
  });

  it('applies custom className', () => {
    render(<Button className="custom-class">Custom</Button>);
    const button = screen.getByText('Custom');
    expect(button).toHaveClass('custom-class');
  });

  it('handles disabled state', () => {
    render(<Button disabled>Disabled</Button>);
    const button = screen.getByText('Disabled');
    expect(button).toBeDisabled();
    expect(button).toHaveClass('disabled:pointer-events-none');
    expect(button).toHaveClass('disabled:opacity-50');
  });

  it('handles click events', async () => {
    const handleClick = vi.fn();
    render(<Button onClick={handleClick}>Click me</Button>);

    await userEvent.click(screen.getByText('Click me'));
    expect(handleClick).toHaveBeenCalledTimes(1);
  });

  it('does not trigger click when disabled', async () => {
    const handleClick = vi.fn();
    render(<Button disabled onClick={handleClick}>Click me</Button>);

    await userEvent.click(screen.getByText('Click me'));
    expect(handleClick).not.toHaveBeenCalled();
  });

  it('forwards additional props', () => {
    render(<Button data-testid="test-button" aria-label="Test">Test</Button>);
    const button = screen.getByText('Test');

    expect(button).toHaveAttribute('data-testid', 'test-button');
    expect(button).toHaveAttribute('aria-label', 'Test');
  });
});
