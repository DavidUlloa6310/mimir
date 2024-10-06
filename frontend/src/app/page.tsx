import LoginCard from '@/components/LoginCard'
import ThemeToggle from '@/components/ThemeToggle'

export default function Home() {
  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center bg-gradient-to-l from-cyan-500 via-teal-600 to-green-700 dark:to-green-950 dark:via-teal-950 dark:from-gray-900">
    <div className="absolute top-4 right-4 z-10">
      <ThemeToggle />
    </div>
    <main className="w-full max-w-md px-2">
      <h1 className="text-7xl font-bold mb-6 text-center">mimir</h1>
      <LoginCard />
    </main>
  </div>
  )
}
