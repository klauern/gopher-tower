import { act, render } from '@testing-library/react';
import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { EventStream, EventStreamData } from '../EventStream';

// Mock EventSource
class MockEventSource {
  onopen: (() => void) | null = null;
  onmessage: ((event: MessageEvent) => void) | null = null;
  onerror: ((error: Event) => void) | null = null;
  readyState = 0;
  url: string;
  withCredentials: boolean;

  constructor(url: string, options?: { withCredentials: boolean }) {
    this.url = url;
    this.withCredentials = options?.withCredentials || false;
  }

  close() {
    // Simulate connection close
    this.readyState = 2;
  }
}

// Replace global EventSource with mock
const originalEventSource = global.EventSource;
global.EventSource = MockEventSource as unknown as typeof EventSource;

describe('EventStream', () => {
  const mockUrl = '/api/events';
  let mockEventSource: MockEventSource;

  beforeEach(() => {
    vi.useFakeTimers();
    mockEventSource = new MockEventSource(mockUrl);
  });

  afterEach(() => {
    vi.restoreAllMocks();
    vi.useRealTimers();
  });

  it('establishes connection on mount', async () => {
    const onConnectionChange = vi.fn();
    render(
      <EventStream url={mockUrl} onConnectionChange={onConnectionChange} />
    );

    // Simulate successful connection
    act(() => {
      mockEventSource.onopen?.();
    });

    expect(onConnectionChange).toHaveBeenCalledWith(true);
  });

  it('handles messages correctly', async () => {
    const onMessage = vi.fn();
    const testData: EventStreamData = {
      type: 'test',
      payload: { message: 'Hello' }
    };

    render(
      <EventStream url={mockUrl} onMessage={onMessage} />
    );

    // Simulate message
    act(() => {
      mockEventSource.onmessage?.({
        data: JSON.stringify(testData)
      } as MessageEvent);
    });

    expect(onMessage).toHaveBeenCalledWith(testData);
  });

  it('handles connection errors', async () => {
    const onError = vi.fn();
    const onConnectionChange = vi.fn();

    render(
      <EventStream
        url={mockUrl}
        onError={onError}
        onConnectionChange={onConnectionChange}
        retryInterval={1000}
      />
    );

    // Simulate error
    act(() => {
      mockEventSource.onerror?.(new Event('error'));
    });

    expect(onError).toHaveBeenCalled();
    expect(onConnectionChange).toHaveBeenCalledWith(false);

    // Fast-forward past retry interval
    act(() => {
      vi.advanceTimersByTime(1000);
    });

    // Verify reconnection attempt
    expect(mockEventSource.readyState).toBe(2); // Closed
  });

  it('handles message parsing errors', async () => {
    const onError = vi.fn();

    render(
      <EventStream url={mockUrl} onError={onError} />
    );

    // Simulate invalid message
    act(() => {
      mockEventSource.onmessage?.({
        data: 'invalid json'
      } as MessageEvent);
    });

    expect(onError).toHaveBeenCalledWith(expect.any(Error));
  });

  it('renders children with connection status', () => {
    const { getByText } = render(
      <EventStream url={mockUrl}>
        {(isConnected) => (
          <div>{isConnected ? 'Connected' : 'Disconnected'}</div>
        )}
      </EventStream>
    );

    expect(getByText('Disconnected')).toBeInTheDocument();

    // Simulate connection
    act(() => {
      mockEventSource.onopen?.();
    });

    expect(getByText('Connected')).toBeInTheDocument();
  });

  it('cleans up on unmount', () => {
    const { unmount } = render(
      <EventStream url={mockUrl} />
    );

    unmount();

    expect(mockEventSource.readyState).toBe(2); // Closed
  });

  it('reconnects when URL changes', () => {
    const { rerender } = render(
      <EventStream url={mockUrl} />
    );

    // Simulate initial connection
    act(() => {
      mockEventSource.onopen?.();
    });

    // Change URL
    rerender(<EventStream url="/api/events/new" />);

    expect(mockEventSource.readyState).toBe(2); // Old connection closed
  });
});

// Restore original EventSource
global.EventSource = originalEventSource;
