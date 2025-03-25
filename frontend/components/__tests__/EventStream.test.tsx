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
    this.readyState = 2;
  }
}

describe('EventStream', () => {
  const mockUrl = '/api/events';
  let mockEventSource: MockEventSource;

  beforeEach(() => {
    vi.useFakeTimers();
    // Mock the global EventSource
    global.EventSource = vi.fn().mockImplementation((url: string, options?: { withCredentials: boolean }) => {
      mockEventSource = new MockEventSource(url, options);
      return mockEventSource;
    }) as unknown as typeof EventSource;
  });

  afterEach(() => {
    vi.clearAllMocks();
    vi.useRealTimers();
  });

  it('establishes connection on mount', () => {
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

  it('handles messages correctly', () => {
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

  it('handles connection errors', () => {
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
    expect(global.EventSource).toHaveBeenCalledTimes(2);
  });

  it('handles message parsing errors', () => {
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

    // Initially disconnected
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

    // Simulate successful connection
    act(() => {
      mockEventSource.onopen?.();
    });

    unmount();

    // Verify cleanup
    expect(mockEventSource.readyState).toBe(2);
  });

  it('reconnects when URL changes', async () => {
    const { rerender } = render(
      <EventStream url={mockUrl} />
    );

    // Wait for initial connection
    await act(async () => {
      mockEventSource.onopen?.();
      await vi.runAllTimersAsync();
    });

    // Verify initial connection
    expect(global.EventSource).toHaveBeenCalledTimes(1);
    expect(global.EventSource).toHaveBeenLastCalledWith(mockUrl, expect.any(Object));

    // Change URL
    await act(async () => {
      rerender(<EventStream url="/api/events/new" />);
      await vi.runAllTimersAsync();
    });

    // Verify new connection was created
    expect(global.EventSource).toHaveBeenCalledTimes(2);
    expect(global.EventSource).toHaveBeenLastCalledWith('/api/events/new', expect.any(Object));
  });
});
