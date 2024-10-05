'use client'

import { useState, useEffect } from 'react'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'

export default function Carousel() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [currentIndex, setCurrentIndex] = useState(0)

  const numberOfCards = 6 

  useEffect(() => {
    // Simulate data fetching with a delay
    const fetchData = async () => {
      setLoading(true)
      await new Promise((resolve) => setTimeout(resolve, 2000)) // Simulate 2-second delay

      const fetchedData = Array.from({ length: numberOfCards }, (_, index) => ({
        id: index + 1,
        title: `Card ${index + 1}`,
        description: `Description ${index + 1}`,
      }))

      setData(fetchedData)
      setLoading(false)
    }

    fetchData()
  }, [numberOfCards]) 

  const itemsToShow = 3

  const next = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex + itemsToShow >= data.length ? prevIndex : prevIndex + itemsToShow
    )
  }

  const prev = () => {
    setCurrentIndex((prevIndex) =>
      prevIndex - itemsToShow <= 0 ? 0 : prevIndex - itemsToShow
    )
  }

  const containerClasses = 'flex space-x-4 overflow-hidden w-full'

  if (loading) {
    return (
      <div className="flex items-center justify-between">
        <Button variant="ghost" size="icon" disabled>
          <ChevronLeft className="w-6 h-6 text-gray-400" />
        </Button>

        <div className={containerClasses}>
          {[1, 2, 3].map((key) => (
            <Skeleton
              key={key}
              className="w-full h-32 bg-white/30 rounded-lg"
            />
          ))}
        </div>

        <Button variant="ghost" size="icon" disabled>
          <ChevronRight className="w-6 h-6 text-gray-400" />
        </Button>
      </div>
    )
  }

  const visibleItems = data.slice(currentIndex, currentIndex + itemsToShow)

  const isPrevDisabled = currentIndex <= 0
  const isNextDisabled = currentIndex + itemsToShow >= data.length

  return (
    <div className="flex items-center justify-between">
      <Button variant="ghost" size="icon" onClick={prev} disabled={isPrevDisabled}>
        <ChevronLeft className={`w-6 h-6 ${isPrevDisabled ? 'text-gray-400' : 'text-white'}`} />
      </Button>

      <div className={containerClasses}>
        {visibleItems.map((item) => (
          <Card
            key={item.id}
            className="w-full p-4 bg-white/70 dark:bg-gray-800/70 backdrop-blur-md shadow-md rounded-lg"
          >
            <h2 className="text-lg font-semibold mb-2">{item.title}</h2>
            <p>{item.description}</p>
          </Card>
        ))}
      </div>

      <Button variant="ghost" size="icon" onClick={next} disabled={isNextDisabled}>
        <ChevronRight className={`w-6 h-6 ${isNextDisabled ? 'text-gray-400' : 'text-white'}`} />
      </Button>
    </div>
  )
}