import LoginCard from '@/components/LoginCard'
import ThemeToggle from '@/components/ThemeToggle'

export default function Home() {
  return (
    <div className="relative min-h-screen flex flex-col items-center justify-center bg-gradient-to-r from-blue-400 to-purple-500 dark:from-blue-900 dark:to-purple-900">
      <div className="absolute top-4 right-4 z-10">
        <ThemeToggle />
      </div>
      <main className="w-full max-w-md px-4">
        <LoginCard />
      </main>
    </div>
  )
}
