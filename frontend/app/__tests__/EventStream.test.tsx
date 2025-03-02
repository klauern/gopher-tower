import { afterEach, beforeEach, describe, expect, it, vi } from 'vitest';
import { EventStream, EventStreamData } from '../components/EventStream';
import { act, render, screen } from '../utils/test-utils';

// Mock EventSource since it's not available in jsdom
class MockEventSource {
  onopen: (() => void) | null = null;
  onmessage: ((event: { data: string }) => void) | null = null;
  onerror: ((error: Event) => void) | null = null;

  constructor(public url: string, public options?: { withCredentials: boolean }) {
    // We'll manually trigger onopen in tests
  }

  // Mock method to simulate receiving a message
  simulateMessage(data: EventStreamData) {
    if (this.onmessage) {
      this.onmessage({ data: JSON.stringify(data) });
    }
  }

  // Mock method to simulate an error
  simulateError() {
    if (this.onerror) {
      this.onerror(new Event('error'));
    }
  }

  close() {
    // Clean up any resources
  }
}

// Create a mock instance we can access in tests
let mockEventSourceInstance: MockEventSource | null = null;

// Mock the EventSource constructor
vi.stubGlobal('EventSource', function(url: string, options?: { withCredentials: boolean }) {
  mockEventSourceInstance = new MockEventSource(url, options);
  return mockEventSourceInstance;
});

describe('EventStream', () => {
  // Mocks
  const onMessage = vi.fn();
  const onConnectionChange = vi.fn();

  // Clear mocks between tests
  beforeEach(() => {
    vi.clearAllMocks();
    vi.useFakeTimers();
    mockEventSourceInstance = null;
  });

  afterEach(() => {
    vi.useRealTimers();
  });

  it('establishes a connection to the provided URL', async () => {
    render(
      <EventStream
        url="/api/events"
        onConnectionChange={onConnectionChange}
      />
    );

    // Trigger the onopen callback
    act(() => {
      mockEventSourceInstance?.onopen?.();
    });

    // Check if onConnectionChange was called with true
    expect(onConnectionChange).toHaveBeenCalledWith(true);
  });

  it('calls onMessage when a message is received', async () => {
    render(
      <EventStream
        url="/api/events"
        onMessage={onMessage}
      />
    );

    // Simulate a message
    const testData: EventStreamData = {
      type: 'test',
      payload: { message: 'Hello, world!' }
    };

    // Trigger the onopen callback
    act(() => {
      mockEventSourceInstance?.onopen?.();
      mockEventSourceInstance?.simulateMessage(testData);
    });

    // Check if onMessage was called with the test data
    expect(onMessage).toHaveBeenCalledWith(testData);
  });

  it('renders children with connection status', async () => {
    render(
      <EventStream url="/api/events">
        {(isConnected) => (
          <div data-testid="connection-status">
            {isConnected ? 'Connected' : 'Disconnected'}
          </div>
        )}
      </EventStream>
    );

    // Initially should be disconnected
    const statusElement = screen.getByTestId('connection-status');
    expect(statusElement.textContent).toBe('Disconnected');

    // Trigger the onopen callback
    act(() => {
      mockEventSourceInstance?.onopen?.();
    });

    // Now should be connected
    const updatedStatusElement = screen.getByTestId('connection-status');
    expect(updatedStatusElement.textContent).toBe('Connected');
  });
});
