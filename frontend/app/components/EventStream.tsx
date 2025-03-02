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
    if (reconnectTimeoutRef.current !== null) {
      window.clearTimeout(reconnectTimeoutRef.current);
      reconnectTimeoutRef.current = null;
    }
    if (eventSourceRef.current) {
      eventSourceRef.current.close();
      eventSourceRef.current = null;
    }
  }, []);

  const connect = useCallback(() => {
    try {
      // Clean up any existing connection
      cleanup();

      console.log('Attempting to connect to SSE endpoint:', url);
      const eventSource = new EventSource(url, { withCredentials: false });
      eventSourceRef.current = eventSource;

      eventSource.onopen = () => {
        console.log('SSE connection established');
        setIsConnected(true);
        onConnectionChange?.(true);
      };

      eventSource.onmessage = (event) => {
        console.log('Raw SSE message received:', event.data);
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
        setIsConnected(false);
        onConnectionChange?.(false);

        // Close the errored connection
        cleanup();

        // Attempt to reconnect after the specified interval
        console.log(`Will attempt to reconnect in ${retryInterval}ms`);
        reconnectTimeoutRef.current = window.setTimeout(() => {
          console.log('Attempting to reconnect...');
          connect();
        }, retryInterval);
      };

      return eventSource;
    } catch (error) {
      console.error('Failed to establish SSE connection:', error);
      onError?.(new Error('Failed to establish SSE connection'));
      return null;
    }
  }, [url, onMessage, onError, onConnectionChange, retryInterval, cleanup]);

  useEffect(() => {
    connect();

    return () => {
      console.log('Cleaning up SSE connection');
      cleanup();
    };
  }, [connect, cleanup]);

  if (children) {
    return <>{children(isConnected)}</>;
  }

  return null;
}
