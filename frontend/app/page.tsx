'use client';

import { useState } from 'react';
import { EventStream, EventStreamData } from './components/EventStream';

export default function Home() {
  const [events, setEvents] = useState<Array<EventStreamData & { type: string }>>([]);

  return (
    <div className="min-h-screen bg-gray-50 dark:bg-gray-900">
      <header className="bg-white dark:bg-gray-800 shadow">
        <div className="max-w-7xl mx-auto py-6 px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold text-gray-900 dark:text-white">
            Event Stream Monitor
          </h1>
        </div>
      </header>

      <main className="max-w-7xl mx-auto py-6 sm:px-6 lg:px-8">
        <div className="px-4 py-6 sm:px-0">
          <div className="border-4 border-dashed border-gray-200 dark:border-gray-700 rounded-lg p-4">
            <EventStream
              url="http://localhost:8080/events"
              onMessage={(data) => {
                console.log('Received event:', data);
                setEvents(prev => [...prev, data as EventStreamData & { type: string }]);
              }}
              onError={(error) => {
                console.error('Stream error:', error);
              }}
              onConnectionChange={(connected) => {
                console.log('Connection state changed:', connected);
              }}
            >
              {(isConnected) => (
                <div className="space-y-4">
                  <div className="flex items-center space-x-2">
                    <div
                      className={`w-3 h-3 rounded-full ${
                        isConnected
                          ? 'bg-green-500 animate-pulse'
                          : 'bg-red-500'
                      }`}
                    />
                    <span className="text-sm font-medium text-gray-700 dark:text-gray-300">
                      {isConnected ? 'Connected' : 'Disconnected'}
                    </span>
                  </div>
                  <div className="bg-white dark:bg-gray-800 shadow rounded-lg p-4">
                    <h2 className="text-lg font-medium text-gray-900 dark:text-white mb-4">
                      Event Log
                    </h2>
                    <div className="h-64 overflow-y-auto font-mono text-sm bg-gray-50 dark:bg-gray-900 p-4 rounded">
                      {events.length === 0 ? (
                        <p className="text-gray-500 dark:text-gray-400">
                          Waiting for events...
                        </p>
                      ) : (
                        <div className="space-y-2">
                          {events.map((event, i) => (
                            <div key={i} className="text-gray-700 dark:text-gray-300">
                              <span className="text-blue-500">{event.type}</span>:{' '}
                              {JSON.stringify(event.payload)}
                            </div>
                          ))}
                        </div>
                      )}
                    </div>
                  </div>
                </div>
              )}
            </EventStream>
          </div>
        </div>
      </main>
    </div>
  );
}
