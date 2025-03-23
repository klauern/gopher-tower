import Link from 'next/link';
import { useRouter } from 'next/router';

export function Navbar() {
  const router = useRouter();

  const isActive = (path: string) => router.pathname === path;

  return (
    <nav className="bg-white dark:bg-gray-800 border-b border-gray-200 dark:border-gray-700">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center">
            <div className="text-xl font-bold text-gray-900 dark:text-white mr-8">
              Gopher Tower
            </div>
            <div className="flex space-x-8">
              <Link
                href="/"
                className={`inline-flex items-center px-3 py-2 text-sm font-medium ${
                  isActive('/')
                    ? 'text-blue-600 dark:text-blue-400'
                    : 'text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100'
                }`}
              >
                Home
              </Link>
              <Link
                href="/jobs"
                className={`inline-flex items-center px-3 py-2 text-sm font-medium ${
                  isActive('/jobs')
                    ? 'text-blue-600 dark:text-blue-400'
                    : 'text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100'
                }`}
              >
                Jobs
              </Link>
            </div>
          </div>
        </div>
      </div>
    </nav>
  );
}
