'use client'

import { useState, useMemo, useEffect } from 'react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Checkbox } from '@/components/ui/checkbox'
import { Card } from '@/components/ui/card'
import ThemeToggle from '@/components/ThemeToggle'
import Carousel from '@/components/carousel'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'
import { Button } from '@/components/ui/button'
import { ExternalLink } from 'lucide-react'
import UserSection from '@/components/UserSection'
import Link from 'next/link'
import { accelerators } from '@/data/accelerators';

// Hardcoded variables for testing
const CATEGORIES = ['Category A', 'Category B', 'Category C']
const TICKETS_PER_CATEGORY = 3
const SINGLE_TICKETS = 10
const PREVIOUS_CHATS_COUNT = 15
const DOCUMENTATION_LINKS_COUNT = 5

export default function Dashboard() {
  const [selectedItems, setSelectedItems] = useState<string[]>([])
  const [isLoading, setIsLoading] = useState(true)

  useEffect(() => {
    const timer = setTimeout(() => {
      setIsLoading(false)
    }, 2000)

    return () => clearTimeout(timer)
  }, [])

  const ticketData = useMemo(() => {
    let id = 1
    const data = CATEGORIES.map(category => ({
      name: category,
      tickets: Array.from({ length: TICKETS_PER_CATEGORY }, () => ({
        id: `ticket-${id}`,
        name: `Ticket ${id++}`
      }))
    }))

    for (let i = 0; i < SINGLE_TICKETS; i++) {
      data.push({ id: `ticket-${id}`, name: `Single Ticket ${id++}` })
    }

    return data
  }, [])

  const handleSelect = (id: string) => {
    setSelectedItems((prev) => {
      const category = ticketData.find(item => item.name === id);
      if (category && category.tickets) {
        const ticketIds = category.tickets.map(ticket => ticket.id);
        if (prev.includes(id)) {
          return prev.filter(itemId => itemId !== id && !ticketIds.includes(itemId));
        } else {
          return [...prev.filter(itemId => !ticketIds.includes(itemId)), id, ...ticketIds];
        }
      } else {
        return prev.includes(id) ? prev.filter((itemId) => itemId !== id) : [...prev, id];
      }
    });
  };

  const previousChats = useMemo(() => {
    return Array.from({ length: PREVIOUS_CHATS_COUNT }, (_, index) => ({
      id: index + 1,
      title: `Chat about ${['React', 'Next.js', 'TypeScript', 'Tailwind CSS'][index % 4]} (${index + 1})`,
      date: new Date(Date.now() - index * 86400000).toISOString().split('T')[0], // Subtracts days
    }))
  }, [])

  const documentationLinks = useMemo(() => {
    const topics = ['React', 'Next.js', 'TypeScript', 'Tailwind CSS', 'Redux']
    return Array.from({ length: DOCUMENTATION_LINKS_COUNT }, (_, index) => ({
      title: `${topics[index % topics.length]} Documentation`,
      url: `https://example.com/${topics[index % topics.length].toLowerCase()}-docs`,
    }))
  }, [])
  return (
    <div className="relative flex h-screen bg-gradient-to-l from-cyan-500 via-teal-600 to-green-700 dark:to-green-950 dark:via-teal-950 dark:from-gray-900">
      <div className="absolute top-4 right-4 z-10">
        <ThemeToggle />
      </div>

      {/* Sidebar */}
      <aside className="h-full w-64 bg-white/30 shadow-black/30 dark:bg-slate-800/80 backdrop-blur-md shadow-md flex flex-col">
        <UserSection />
        <Separator className="mt-10  bg-gray-400 shadow-black dark:bg-gray-600" />
        <ScrollArea className="flex-1">
          <div className="p-4 space-y-2">
            {ticketData.map((item, index) => (
              <div key={item.id || `category-${index}`}>
                {item.tickets ? (
                  <>
                    <div className="flex items-center justify-between p-2 hover:bg-gray-200/50 dark:hover:bg-gray-700 rounded-md">
                      <h3 className="font-semibold text-md text-gray-800 dark:text-gray-200">
                        {item.name}
                      </h3>
                      <Checkbox
                        checked={selectedItems.includes(item.name)}
                        onCheckedChange={() => handleSelect(item.name)}
                        id={item.name}
                      />
                    </div>
                    <Separator className="my-1 bg-gray-400 dark:bg-gray-600 shadow-md" />
                    <div className="ml-4 space-y-1">
                      {item.tickets.map((ticket) => (
                        <div key={ticket.id}>
                          <div className="flex items-center justify-between p-2 hover:bg-gray-200/50 dark:hover:bg-gray-700 rounded-md">
                            <label
                              htmlFor={ticket.id}
                              className="text-gray-600 dark:text-gray-200"
                            >
                              {ticket.name}
                            </label>
                            <Checkbox
                              checked={selectedItems.includes(ticket.id)}
                              onCheckedChange={() => handleSelect(ticket.id)}
                              id={ticket.id}
                            />
                          </div>
                          <Separator className="my-1 bg-gray-400 dark:bg-gray-600" />
                        </div>
                      ))}
                    </div>
                  </>
                ) : (
                  <div>
                    <div className="flex items-center justify-between p-2 hover:bg-gray-200/50 dark:hover:bg-gray-700 rounded-md">
                      <label
                        htmlFor={item.id}
                        className="text-gray-700 dark:text-gray-200"
                      >
                        {item.name}
                      </label>
                      <Checkbox
                        checked={selectedItems.includes(item.id)}
                        onCheckedChange={() => handleSelect(item.id)}
                        id={item.id}
                      />
                    </div>
                    <Separator className="my-2 bg-gray-400 dark:bg-gray-600" />
                  </div>
                )}
              </div>
            ))}
          </div>
        </ScrollArea>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-4 overflow-auto">
        <div className="grid grid-cols-2 grid-rows-6 gap-4 h-full">
          {/* Dashboard Header - 1 row, 2 columns */}
          <Card className="col-span-2 bg-white/30 dark:bg-gray-800/80 backdrop-blur-md shadow-md flex items-center justify-center">
            <h1 className="text-7xl font-black italic text-gray-200  dark:text-gray-400 p-4">Dashboard</h1>
          </Card>

          {/* Carousel - 2 rows, 2 columns */}
          <Card className="col-span-2 row-span-2 bg-white/30 dark:bg-slate-800/80 backdrop-blur-md shadow-md p-4 flex items-center justify-center">
            <div className="w-full h-full">
              <Carousel />
            </div>
          </Card>

          {/* Previous Chats - 3 rows, 1 column */}
          <Card className="row-span-3 p-4 bg-white/30 dark:bg-slate-800/80 backdrop-blur-md shadow-md flex flex-col">
            <h2 className="text-2xl font-semibold mb-2 text-gray-800 dark:text-gray-200">Previous Chats</h2>
            <Separator className="mb-4 bg-gray-400 dark:bg-gray-600 shadow-emerald-50" />
            {isLoading ? (
              <div className="flex-1 flex flex-col justify-around">
                <Skeleton className="w-full h-1/5" />
                <Skeleton className="w-full h-1/5" />
                <Skeleton className="w-full h-1/5" />
              </div>
            ) : (
              <ScrollArea className="flex-1">
                <div className="flex flex-col gap-4">
                  {previousChats.map((chat) => (
                    <Link key={chat.id} href="/chatpage">
                      <Card className="p-3 bg-white/50 dark:bg-gray-700/50 hover:bg-white/60 dark:hover:bg-gray-700/60 transition-colors cursor-pointer  ">
                        <h3 className="font-semibold text-gray-800 dark:text-gray-200">{chat.title}</h3>
                        <p className="text-sm text-gray-600 dark:text-gray-400">{chat.date}</p>
                      </Card>
                    </Link>
                  ))}
                </div>
              </ScrollArea>
            )}
          </Card>

          {/* Documentation - 3 rows, 1 column */}
          <Card className="row-span-3 p-4 bg-white/30 dark:bg-slate-800/80 backdrop-blur-md shadow-md flex flex-col">
            <h2 className="text-2xl font-semibold mb-2 text-gray-800 dark:text-gray-200">Documentation</h2>
            <Separator className="mb-4 bg-gray-400 dark:bg-gray-600 shadow-emerald-50" />
            {isLoading ? (
              <div className="flex-1 flex flex-col justify-around">
                <Skeleton className="w-full h-1/5" />
                <Skeleton className="w-full h-1/5" />
                <Skeleton className="w-full h-1/5" />
              </div>
            ) : (
              <ScrollArea className="flex-1 h-64">
                <div className="p-4 space-y-2">
                  {accelerators.map((accelerator, index) => (
                    <div key={index}>
                      <div className="flex items-center justify-between p-2 hover:bg-gray-200/50 dark:hover:bg-gray-700 rounded-md">
                        <Link
                          href={accelerator.url}
                          target="_blank"
                          rel="noopener noreferrer"
                          className="text-blue-600 dark:text-blue-400 hover:underline"
                        >
                          {accelerator.title}
                        </Link>
                      </div>
                      <Separator className="my-1 bg-gray-400 dark:bg-gray-600 shadow-md" />
                    </div>
                  ))}
                </div>
              </ScrollArea>
            )}
          </Card>
        </div>
      </main>
    </div>
  )
}