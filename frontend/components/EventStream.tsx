'use client';

import { ReactNode, useCallback, useEffect, useRef, useState } from 'react';

export type EventStreamData = {
  type: string;
  payload: Record<string, unknown>;
};

interface EventStreamProps {
  url: string;
  onMessage?: (data: EventStreamData) => void;
  onError?: (error: Error) => void;
  onConnectionChange?: (isConnected: boolean) => void;
  retryInterval?: number;
  children?: (isConnected: boolean) => ReactNode;
}

export function EventStream({
  url,
  onMessage,
  onError,
  onConnectionChange,
  retryInterval = 5000,
  children,
}: EventStreamProps) {
  const [isConnected, setIsConnected] = useState(false);
  const eventSourceRef = useRef<EventSource | null>(null);
  const reconnectTimeoutRef = useRef<number | null>(null);

  const cleanup = useCallback(() => {
    console.log('Cleaning up EventStream resources');
    if (reconnectTimeoutRef.current !== null) {
      console.log('Clearing reconnect timeout');
      window.clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (eventSourceRef.current) {
      console.log('Closing existing EventSource connection');
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
  }, []);

  const connect = useCallback(() => {
    try {
      // Clean up any existing connection
      cleanup();

      console.log('Creating new EventSource connection to:', url);
      const eventSource = new EventSource(url, { withCredentials: false });
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        console.log('SSE connection established successfully');
        setIsConnected(true);
        onConnectionChange?.(true);
      };

      eventSource.onmessage = (event) => {
        console.log('SSE message received:', {
          type: event.type,
          data: event.data,
          lastEventId: event.lastEventId,
        });
        try {
          const data = JSON.parse(event.data) as EventStreamData;
          console.log('Parsed SSE message:', data);
          onMessage?.(data);
        } catch (error) {
          console.error('Failed to parse SSE message:', error);
          console.error('Raw message that failed:', event.data);
          onError?.(new Error('Failed to parse SSE message'));
        }
      };

      eventSource.onerror = (error) => {
        console.error('SSE connection error:', error);
        console.log('EventSource readyState:', eventSource.readyState);
        setIsConnected(false);
        onConnectionChange?.(false);
        onError?.(new Error('SSE connection error'));

        // Close the errored connection
        cleanup();

        // Attempt to reconnect after the specified interval
        console.log(`Scheduling reconnect attempt in ${retryInterval}ms`);
        reconnectTimeoutRef.current = window.setTimeout(() => {
          console.log('Attempting to reconnect to SSE...');
          connect();
        }, retryInterval);
      };

      return eventSource;
    } catch (error) {
      console.error('Failed to create EventSource:', error);
      onError?.(new Error('Failed to establish SSE connection'));
      return null;
    }
  }, [url, onMessage, onError, onConnectionChange, retryInterval, cleanup]);

  useEffect(() => {
    console.log('EventStream component mounted, initializing connection');
    connect();

    return () => {
      console.log('EventStream component unmounting');
      cleanup();
    };
  }, [connect, cleanup]);

  // Add a reconnection effect when the URL changes
  useEffect(() => {
    if (isConnected) {
      console.log('URL changed, reconnecting...');
      connect();
    }
  }, [url, connect]);

  if (children) {
    return <>{children(isConnected)}</>;
  }

  return null;
}
