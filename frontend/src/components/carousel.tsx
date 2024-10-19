"use client";

import { useState, useEffect } from "react";
import { Button } from "@/components/ui/button";
import {
  ChevronLeft,
  ChevronRight,
  Rocket,
  Wrench,
  Glasses,
  Lightbulb,
} from "lucide-react";
import { motion } from "framer-motion";
import {
  Tooltip,
  TooltipContent,
  TooltipProvider,
  TooltipTrigger,
} from "@/components/ui/tooltip";
import Link from "next/link";
import { accelerators } from "@/data/accelerators";

interface CarouselCardProps {
  name: string;
  description: string;
  url: string;
  category: "Architecture" | "Strategy" | "Technical";
  acceleratorId: number;
}

export default function Carousel() {
  const [data, setData] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);
  const [currentIndex, setCurrentIndex] = useState(0);
  const [key, setKey] = useState(0);

  const itemsToShow = 3;

  useEffect(() => {
    setLoading(true);
    // Simulate a delay to mimic async data fetching
    setTimeout(() => {
      setData(accelerators);
      setLoading(false);
    }, 500); // Adjust delay as needed
  }, []);

  const next = () => {
    setCurrentIndex((prevIndex) => {
      const maxIndex = Math.max(0, data.length - itemsToShow);
      const newIndex =
        prevIndex + itemsToShow >= data.length
          ? maxIndex
          : prevIndex + itemsToShow;
      setKey((prev) => prev + 1);
      return newIndex;
    });
  };

  const prev = () => {
    setCurrentIndex((prevIndex) => {
      const newIndex =
        prevIndex - itemsToShow <= 0 ? 0 : prevIndex - itemsToShow;
      setKey((prev) => prev + 1);
      return newIndex;
    });
  };

  const containerVariants = {
    hidden: { opacity: 0 },
    visible: {
      opacity: 1,
      transition: { staggerChildren: 0.1 },
    },
  };

  const cardVariants = {
    hidden: { opacity: 0, y: 20 },
    visible: {
      opacity: 1,
      y: 0,
      transition: { duration: 0.5 },
    },
  };

  if (loading) {
    return (
      <div className="flex items-center justify-between w-full h-full gap-2">
        <Button variant="outline" size="icon" disabled>
          <ChevronLeft className="w-6 h-6 text-gray-400" />
        </Button>

        <div className="flex-1 grid grid-cols-3 gap-4">
          {[1, 2, 3].map((key) => (
            <div
              key={key}
              className="w-full h-64 bg-gray-200 dark:bg-gray-700 rounded-lg animate-pulse"
            />
          ))}
        </div>

        <Button variant="outline" size="icon" disabled>
          <ChevronRight className="w-6 h-6 text-gray-400" />
        </Button>
      </div>
    );
  }

  const visibleItems = data.slice(currentIndex, currentIndex + itemsToShow);

  const isPrevDisabled = currentIndex <= 0;
  const isNextDisabled = currentIndex + itemsToShow >= data.length;
  
  function CarouselCard({
    name,
    description,
    url,
    category,
    acceleratorId
  }: CarouselCardProps) {
    const categoryIcons = {
      Architecture: Rocket,
      Strategy: Wrench,
      Technical: Glasses,
      // Add more categories and corresponding icons as needed
    };

    const IconComponent = categoryIcons[category] || Lightbulb;

    const handleClick = async () => {
      // Define your credentials
      const username = "admin"; // Replace with your username
      const password = "r8RGnqYX=%m0";

      // Create the thread
      const response = await fetch(`${process.env.NEXT_PUBLIC_BACKEND_IP}/chat`, {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
          Authorization: "Basic " + btoa(`${username}:${password}`),
        },
        body: JSON.stringify({
          instanceId: "dev274800",
          createThread: true,
          acceleratorId: acceleratorId,
        }),
      });

      if (!response.ok) {
        console.error("Failed to create thread:", response.statusText);
        return;
      }

      const { threadId } = await response.json();

      window.location.href = `/chatpage?threadId=${threadId}&acceleratorId=${acceleratorId}`;
    };

    return (
      <div
        className="carousel-card bg-white/60 dark:bg-gray-700/60 rounded-lg  max-h-[250px] shadow-md p-6 w-full h-full flex flex-col items-center justify-center relative cursor-pointer hover:bg-white/70 dark:hover:bg-gray-700/70 transition-colors"
        onClick={handleClick}
      >
        <div className="w-full h-full flex flex-col items-center justify-center">
          <IconComponent className="w-12 h-12 text-blue-500 mb-4" />
          <div className="text-center">
            <h3 className="text-xl font-semibold mb-2">{name}</h3>
            <p className="text-gray-600 dark:text-gray-300">{description}</p>
          </div>
        </div>
        <TooltipProvider>
          <Tooltip>
            <TooltipTrigger asChild>
              <div className="absolute bottom-0 left-0 z-100 bg-gray-400/30 dark:bg-gray-600 p-2 rounded-tr-lg rounded-bl-lg">
                <Lightbulb
                  size={20}
                  className="text-gray-600 dark:text-gray-300"
                />
              </div>
            </TooltipTrigger>
            <TooltipContent>
              <p>Click to learn more</p>
            </TooltipContent>
          </Tooltip>
        </TooltipProvider>
      </div>
    );
  }

  return (
    <div className="flex items-center justify-between w-full h-full gap-2">
      <Button
        variant="outline"
        size="icon"
        onClick={prev}
        disabled={isPrevDisabled}
      >
        <ChevronLeft
          className={`w-6 h-6 ${
            isPrevDisabled ? "text-gray-600" : "text-gray-400"
          }`}
        />
      </Button>

      <motion.div
        className="flex-1 grid grid-cols-3 gap-4"
        variants={containerVariants}
        initial="hidden"
        animate="visible"
        key={key}
      >
        {visibleItems.map((item) => (
          <motion.div
            key={item.id}
            variants={cardVariants}
            className="w-full h-full"
          >
            <CarouselCard
              name={item.name}
              acceleratorId={item._additional.id}
              description={item.description}
              url={item.url}
              category={item.category}
            />
          </motion.div>
        ))}
      </motion.div>

      <Button
        variant="outline"
        size="icon"
        onClick={next}
        disabled={isNextDisabled}
      >
        <ChevronRight
          className={`w-6 h-6 ${
            isNextDisabled ? "text-gray-600" : "text-gray-400"
          }`}
        />
      </Button>
    </div>
  );
}
