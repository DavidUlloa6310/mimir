'use client'

import { useState, useEffect } from 'react'
import { Card } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { ChevronLeft, ChevronRight } from 'lucide-react'
import { Skeleton } from '@/components/ui/skeleton'
import { motion } from 'framer-motion'

export default function Carousel() {
  const [data, setData] = useState<any[]>([])
  const [loading, setLoading] = useState(true)
  const [currentIndex, setCurrentIndex] = useState(0)
  const [key, setKey] = useState(0)

  const numberOfCards = 6 

  useEffect(() => {
    // Simulate data fetching with a delay
    const fetchData = async () => {
      setLoading(true)
      await new Promise((resolve) => setTimeout(resolve, 2000)) 

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
    setCurrentIndex((prevIndex) => {
      const newIndex = prevIndex + itemsToShow >= data.length ? prevIndex : prevIndex + itemsToShow
      setKey(prev => prev + 1)
      return newIndex
    })
  }

  const prev = () => {
    setCurrentIndex((prevIndex) => {
      const newIndex = prevIndex - itemsToShow <= 0 ? 0 : prevIndex - itemsToShow
      setKey(prev => prev + 1)
      return newIndex
    })
  }

  const containerClasses = 'flex space-x-4 w-full h-full'

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: { 
      opacity: 1,
      transition: { staggerChildren: 0.1 }
    }
  }

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: { 
      opacity: 1, 
      y: 0,
      transition: { duration: 0.5 }
    }
  }

  if (loading) {
    return (
      <div className="flex items-center justify-between w-full h-full gap-2">
        <Button variant="outline" size="icon" disabled>
          <ChevronLeft className="w-6 h-6 text-gray-400" />
        </Button>

        <div className={containerClasses}>
          {[1, 2, 3].map((key) => (
            <Skeleton
              key={key}
              className="w-full h-full rounded-lg"
            />
          ))}
        </div>

        <Button variant="outline" size="icon" disabled>
          <ChevronRight className="w-6 h-6 text-gray-400" />
        </Button>
      </div>
    )
  }

  const visibleItems = data.slice(currentIndex, currentIndex + itemsToShow)

  const isPrevDisabled = currentIndex <= 0
  const isNextDisabled = currentIndex + itemsToShow >= data.length

  return (
    <div className="flex items-center justify-between w-full h-full gap-2">
      <Button variant="outline" size="icon" onClick={prev} disabled={isPrevDisabled}>
        <ChevronLeft className={`w-6 h-6 ${isPrevDisabled ? 'text-gray-400' : 'text-gray-700'}`} />
      </Button>

      <motion.div 
        className={containerClasses}
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        key={key}
      >
        {visibleItems.map((item) => (
          <motion.div
            key={`${item.id}-${key}`}
            variants={cardVariants}
            className="w-full h-full" // Add this to maintain original size
          >
            <Card
              className="w-full h-full p-4 bg-white/20 dark:bg-gray-700/80 backdrop-blur-md shadow-md rounded-lg flex flex-col justify-between"
            >
              <h2 className="text-lg font-semibold mb-2">{item.title}</h2>
              <p>{item.description}</p>
            </Card>
          </motion.div>
        ))}
      </motion.div>

      <Button variant="outline" size="icon" onClick={next} disabled={isNextDisabled}>
        <ChevronRight className={`w-6 h-6 ${isNextDisabled ? 'text-gray-400' : 'text-gray-700'}`} />
      </Button>
    </div>
  )
}