'use client'

import {
  NavigationMenu,
  NavigationMenuItem,
  NavigationMenuLink,
  NavigationMenuList,
  navigationMenuTriggerStyle
} from "@/components/ui/navigation-menu"
import { cn } from "@/lib/utils"
import Link from 'next/link'
import { usePathname } from 'next/navigation'
import { ThemeToggle } from './theme-toggle'

export function Navbar() {
  const pathname = usePathname()

  return (
    <div className="border-b border-border bg-background">
      <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
        <div className="flex h-16 items-center justify-between">
          <div className="flex items-center gap-8">
            <Link href="/" className="text-xl font-bold">
              Gopher Tower
            </Link>
            <NavigationMenu>
              <NavigationMenuList>
                <NavigationMenuItem>
                  <Link href="/" legacyBehavior passHref>
                    <NavigationMenuLink className={cn(
                      navigationMenuTriggerStyle(),
                      "text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100",
                      pathname === "/" && "!text-blue-600 dark:!text-blue-400"
                    )}>
                      Home
                    </NavigationMenuLink>
                  </Link>
                </NavigationMenuItem>
                <NavigationMenuItem>
                  <Link href="/jobs" legacyBehavior passHref>
                    <NavigationMenuLink className={cn(
                      navigationMenuTriggerStyle(),
                      "text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100",
                      pathname === "/jobs" && "!text-blue-600 dark:!text-blue-400"
                    )}>
                      Jobs
                    </NavigationMenuLink>
                  </Link>
                </NavigationMenuItem>
                <NavigationMenuItem>
                  <Link href="/events" legacyBehavior passHref>
                    <NavigationMenuLink className={cn(
                      navigationMenuTriggerStyle(),
                      "text-gray-500 dark:text-gray-300 hover:text-gray-700 dark:hover:text-gray-100",
                      pathname === "/events" && "!text-blue-600 dark:!text-blue-400"
                    )}>
                      Events
                    </NavigationMenuLink>
                  </Link>
                </NavigationMenuItem>
              </NavigationMenuList>
            </NavigationMenu>
          </div>
          <div>
            <ThemeToggle />
          </div>
        </div>
      </div>
    </div>
  )
}
