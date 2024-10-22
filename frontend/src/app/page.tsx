"use client";

import LoginCard from '@/components/LoginCard'
import ThemeToggle from '@/components/ThemeToggle'
import { useEffect, useRef } from 'react'

export default function Home() {
  const eyeRef = useRef<SVGCircleElement | null>(null);
  const pupilRef = useRef<SVGCircleElement | null>(null);

  useEffect(() => {
    const handleMouseMove = (e: MouseEvent) => {
      if (eyeRef.current && pupilRef.current) {
        const eyeRect = eyeRef.current.getBoundingClientRect();
        const eyeCenterX = eyeRect.left + eyeRect.width / 2;
        const eyeCenterY = eyeRect.top + eyeRect.height / 2;

        const angle = Math.atan2(e.clientY - eyeCenterY, e.clientX - eyeCenterX);
        const distance = Math.min(eyeRect.width / 4, Math.hypot(e.clientX - eyeCenterX, e.clientY - eyeCenterY) / 5);

        const newX = Math.cos(angle) * distance;
        const newY = Math.sin(angle) * distance;

        eyeRef.current.setAttribute("cx", `${12 + newX}`);
        pupilRef.current.setAttribute("cx", `${13 + newX}`);
        eyeRef.current.setAttribute("cy", `${12 + newY}`);
        pupilRef.current.setAttribute("cy", `${11 + newY}`);
      }
    };

    window.addEventListener('mousemove', handleMouseMove);

    return () => {
      window.removeEventListener('mousemove', handleMouseMove);
    };
  }, []);

  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center bg-gradient-to-l from-cyan-500 via-teal-600 to-green-700 dark:to-green-950 dark:via-teal-950 dark:from-gray-900">
      <div className="absolute top-4 right-4 z-10">
        <ThemeToggle />
      </div>
      <main className="w-full max-w-md px-2">
        <h1 className="text-7xl font-bold mb-6 text-center">mimir</h1>
        <LoginCard eyeRef={eyeRef} pupilRef={pupilRef} />
      </main>
    </div>
  )
}
