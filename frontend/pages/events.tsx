import { EventStream, EventStreamData } from '@/components/EventStream';
import { Navbar } from '@/components/Navbar';
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';
import { getApiUrl } from '@/config';
import { useState } from 'react';

export default function Events() {
  const [events, setEvents] = useState<Array<EventStreamData & { type: string }>>([]);

  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold mb-8">
            Event Stream Monitor
          </h1>

          <Card>
            <CardHeader>
              <CardTitle>Live Events</CardTitle>
            </CardHeader>
            <CardContent>
              <EventStream
                url={getApiUrl('events')}
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
                        className={`w-3 h-3 rounded-full ${isConnected
                          ? 'bg-green-500 animate-pulse'
                          : 'bg-red-500'
                          }`}
                      />
                      <span className="text-sm font-medium text-muted-foreground">
                        {isConnected ? 'Connected' : 'Disconnected'}
                      </span>
                    </div>
                    <div className="rounded-lg">
                      <h2 className="text-lg font-medium mb-4">
                        Event Log
                      </h2>
                      <div className="h-64 overflow-y-auto font-mono text-sm bg-muted p-4 rounded">
                        {events.length === 0 ? (
                          <p className="text-muted-foreground">
                            Waiting for events...
                          </p>
                        ) : (
                          <div className="space-y-2">
                            {events.map((event, i) => (
                              <div key={i} className="text-foreground">
                                <span className="text-primary">{event.type}</span>:{' '}
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
            </CardContent>
          </Card>
        </div>
      </div>
    </div>
  );
}
