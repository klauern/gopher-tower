import { Navbar } from '@/components/Navbar';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import Link from 'next/link';

export default function Home() {
  return (
    <div className="min-h-screen bg-background">
      <Navbar />
      <div className="py-6">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <h1 className="text-3xl font-bold mb-8">
            Dashboard
          </h1>

          <div className="grid gap-4">
            <Card>
              <CardHeader>
                <CardTitle>Welcome to Gopher Tower</CardTitle>
                <CardDescription>
                  Your system monitoring dashboard is coming soon
                </CardDescription>
              </CardHeader>
              <CardContent>
                <p className="text-muted-foreground">
                  The dashboard will provide an overview of system status, key metrics, and quick access to important features.
                  For now, you can:
                </p>
                <ul className="list-disc list-inside mt-4 space-y-2 text-muted-foreground">
                  <li>View real-time events in the <Link href="/events" className="text-primary hover:underline">Event Stream</Link></li>
                  <li>Manage your jobs in the <Link href="/jobs" className="text-primary hover:underline">Jobs</Link> section</li>
                </ul>
              </CardContent>
            </Card>
          </div>
        </div>
      </div>
    </div>
  );
}
