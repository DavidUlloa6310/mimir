'use client'

import { useState, useMemo } from 'react'
import { ScrollArea } from '@/components/ui/scroll-area'
import { Checkbox } from '@/components/ui/checkbox'
import { Card } from '@/components/ui/card'
import ThemeToggle from '@/components/ThemeToggle'
import Carousel from '@/components/carousel'
import { Skeleton } from '@/components/ui/skeleton'
import { Separator } from '@/components/ui/separator'

// Hardcoded variables for testing
const CATEGORIES = ['Category A', 'Category B', 'Category C']
const TICKETS_PER_CATEGORY = 3
const SINGLE_TICKETS = 10

export default function Dashboard() {
  const [selectedItems, setSelectedItems] = useState<string[]>([])

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
    setSelectedItems((prev) =>
      prev.includes(id) ? prev.filter((itemId) => itemId !== id) : [...prev, id]
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
                    <Separator className="my-1" />
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
                          <Separator className="my-1 opacity-50" />
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
                    <Separator className="my-2" />
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
          <Card className="col-span-2 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md flex items-center justify-center">
            <h1 className="text-7xl font-black italic text-gray-200  dark:text-gray-400 p-4">Dashboard</h1>
          </Card>

          {/* Carousel - 2 rows, 2 columns */}
          <Card className="col-span-2 row-span-2 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md p-4 flex items-center justify-center">
            <div className="w-full h-full">
              <Carousel />
            </div>
          </Card>

          {/* Placeholder Skeletons - 3 rows, 1 column each */}
          <Card className="row-span-3 p-4 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
            <Skeleton className="w-full h-full" />
          </Card>
          <Card className="row-span-3 p-4 bg-white/30 dark:bg-gray-800/30 backdrop-blur-md shadow-md">
            <Skeleton className="w-full h-full" />
          </Card>
        </div>
      </main>
    </div>
  )
}