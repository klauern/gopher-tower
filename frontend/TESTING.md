# Testing in the Next.js Frontend

This project uses Vitest with React Testing Library for testing React components in the Next.js application.

## Testing Stack

- **Vitest**: A fast Vite-native testing framework
- **React Testing Library**: A library for testing React components in a user-centric way
- **JSDOM**: A JavaScript implementation of the DOM for testing in Node.js

## Running Tests

You can run tests using either the Taskfile commands (recommended) or directly with Bun:

### Using Taskfile (Recommended)

```bash
# Run all tests (frontend and backend)
task test

# Run only frontend tests
task frontend:test

# Run frontend tests in watch mode
task frontend:test:watch

# Run frontend tests with UI
task frontend:test:ui

# Run frontend tests with coverage report
task frontend:test:coverage
```

### Using Bun Directly

```bash
# Navigate to the frontend directory
cd frontend

# Run all tests
bun run test

# Run tests in watch mode
bun run test:watch

# Run tests with UI
bun run test:ui

# Run tests with coverage report
bun run test:coverage
```

## Test File Structure

- Test files should be placed in the `app/__tests__` directory
- Test files should be named with `.test.tsx` or `.test.ts` extension
- You can also place test files next to the components they test with the same naming convention

## Writing Tests

### Basic Component Test

```tsx
import { describe, it, expect } from 'vitest';
import { render, screen } from '../utils/test-utils';
import MyComponent from '../components/MyComponent';

describe('MyComponent', () => {
  it('renders correctly', () => {
    render(<MyComponent />);
    expect(screen.getByText('Expected Text')).toBeDefined();
  });
});
```

### Testing User Interactions

```tsx
import { describe, it, expect } from 'vitest';
import { render, screen } from '../utils/test-utils';
import userEvent from '@testing-library/user-event';
import Counter from '../components/Counter';

describe('Counter', () => {
  it('increments count when button is clicked', async () => {
    const user = userEvent.setup();
    render(<Counter />);

    const button = screen.getByRole('button', { name: /increment/i });
    await user.click(button);

    expect(screen.getByText('Count: 1')).toBeDefined();
  });
});
```

### Testing Asynchronous Code

```tsx
import { describe, it, expect } from 'vitest';
import { render, screen, waitFor } from '../utils/test-utils';
import AsyncComponent from '../components/AsyncComponent';

describe('AsyncComponent', () => {
  it('loads data and displays it', async () => {
    render(<AsyncComponent />);

    // Initially shows loading state
    expect(screen.getByText('Loading...')).toBeDefined();

    // Wait for data to load
    await waitFor(() => {
      expect(screen.getByText('Data loaded!')).toBeDefined();
    });
  });
});
```

## Mocking

### Mocking Components

```tsx
import { describe, it, expect, vi } from 'vitest';
import { render, screen } from '../utils/test-utils';
import ParentComponent from '../components/ParentComponent';

// Mock a child component
vi.mock('../components/ChildComponent', () => ({
  ChildComponent: () => <div data-testid="mocked-child">Mocked Child</div>
}));

describe('ParentComponent', () => {
  it('renders with mocked child', () => {
    render(<ParentComponent />);
    expect(screen.getByTestId('mocked-child')).toBeDefined();
  });
});
```

### Mocking API Calls

```tsx
import { describe, it, expect, vi, beforeEach } from 'vitest';
import { render, screen, waitFor } from '../utils/test-utils';
import DataComponent from '../components/DataComponent';

// Mock fetch
global.fetch = vi.fn();

describe('DataComponent', () => {
  beforeEach(() => {
    vi.resetAllMocks();
  });

  it('fetches and displays data', async () => {
    // Mock successful response
    global.fetch.mockResolvedValueOnce({
      ok: true,
      json: async () => ({ data: 'Test Data' }),
    });

    render(<DataComponent />);

    await waitFor(() => {
      expect(screen.getByText('Test Data')).toBeDefined();
    });

    expect(global.fetch).toHaveBeenCalledWith('/api/data');
  });
});
```

## Best Practices

1. **Test behavior, not implementation**: Focus on what the component does, not how it does it.
2. **Use user-centric queries**: Prefer queries like `getByRole`, `getByLabelText`, and `getByText` over `getByTestId`.
3. **Mock external dependencies**: Mock API calls, third-party components, and browser APIs.
4. **Keep tests isolated**: Each test should be independent of others.
5. **Test edge cases**: Test loading states, error states, and boundary conditions.
6. **Use act() for state updates**: Wrap code that causes React state updates in `act()`.

## Resources

- [Vitest Documentation](https://vitest.dev/)
- [React Testing Library Documentation](https://testing-library.com/docs/react-testing-library/intro/)
- [Testing Library Queries Cheatsheet](https://testing-library.com/docs/react-testing-library/cheatsheet/)
