'use client'

import { useState } from 'react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Checkbox } from '@/components/ui/checkbox'
import { Card } from '@/components/ui/card'
import ThemeToggle from '@/components/ThemeToggle'
import Carousel from '@/components/carousel'
import { Skeleton } from '@/components/ui/skeleton'

export default function Dashboard() {
  const [selectedTickets, setSelectedTickets] = useState<string[]>([])

  const numberOfTickets = 40

  const tickets = Array.from({ length: numberOfTickets }, (_, index) => ({
    id: `ticket${index + 1}`,
    name: `Sample Ticket ${index + 1}`,
  }))

  const handleSelect = (id: string) => {
    setSelectedTickets((prev) =>
      prev.includes(id) ? prev.filter((ticketId) => ticketId !== id) : [...prev, id]
    )
  }

  return (
    <div className="relative flex h-screen bg-gradient-to-r from-blue-400 to-purple-500 dark:from-blue-900 dark:to-purple-900">
      <div className="absolute top-4 right-4 z-10">
        <ThemeToggle />
      </div>

      {/* Sidebar */}
      <aside className="h-full w-64 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
        <ScrollArea className="h-full">
          <div className="p-4 space-y-2">
            {tickets.map((ticket) => (
              <div
                key={ticket.id}
                className="flex items-center p-2 hover:bg-gray-100 dark:hover:bg-gray-700 rounded-md"
              >
                <Checkbox
                  checked={selectedTickets.includes(ticket.id)}
                  onCheckedChange={() => handleSelect(ticket.id)}
                  id={ticket.id}
                />
                <label
                  htmlFor={ticket.id}
                  className="ml-2 text-gray-800 dark:text-gray-200"
                >
                  {ticket.name}
                </label>
              </div>
            ))}
          </div>
        </ScrollArea>
      </aside>

      {/* Main Content */}
      <main className="flex-1 p-4 overflow-auto">
        <div className="grid grid-cols-1 gap-4 h-full">
          <div className="grid grid-rows-3">
            {/* Dashboard Header */}
            <Card className="bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
              <h1 className="text-2xl font-bold text-white p-4">Dashboard</h1>
            </Card>
          </div>
          <div className="grgrid-rows-3">
            {/* Carousel */}
            <Card className="bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md p-4 flex items-center justify-center">
              <Carousel />
            </Card>
          </div>

          {/* Placeholder Skeletons */}
          <div className="grid grid-cols-2 gap-4">
            <Card className="p-4 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
              <Skeleton className="w-full h-full" />
            </Card>
            <Card className="p-4 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
              <Skeleton className="w-full h-full" />
            </Card>
          </div>
        </div>
      </main>
    </div>
  )
}