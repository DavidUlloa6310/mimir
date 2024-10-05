'use client'

import { useState } from 'react'
import { useRouter } from 'next/navigation'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'
import { Eye, EyeOff } from 'lucide-react'

export default function LoginCard() {
  const router = useRouter()
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [instanceId, setInstanceId] = useState('')
  const [showPassword, setShowPassword] = useState(false)

  const isFormValid = username && password && instanceId

  const handlePasswordToggle = () => {
    setShowPassword((prev) => !prev)
  }

  const handleLogin = (e: React.FormEvent) => {
    e.preventDefault()
    // Temporary code for prototyping: Navigate to the dashboard page
    router.push('/dashboard')
  }

  const inputClasses = 'h-12 bg-white/50 dark:bg-gray-700/50'

  return (
    <Card className="w-full backdrop-blur-md bg-white/30 dark:bg-gray-800/30 shadow-xl">
      <CardHeader>
        {/* <CardTitle className="text-2xl font-bold text-center">Login</CardTitle> */}
        <CardDescription className="text-center">Enter your ServiceNow credentials</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4" onSubmit={handleLogin}>
          {/* Instance ID Input */}
          <div>
            <Input
              type="text"
              placeholder="Instance ID"
              value={instanceId}
              onChange={(e) => setInstanceId(e.target.value)}
              className={inputClasses}
            />
          </div>
          <div className="border-b border-gray-700 dark:border-gray-200"></div>
          {/* Username Input */}
          <div>
            <Input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className={inputClasses}
            />
          </div>
          {/* Password Input with Fade Effect */}
          <div className="relative h-12">
            {/* Password Input (password type) */}
            <Input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={`${inputClasses} pr-10 w-full absolute inset-0 transition-opacity duration-300 ease-in-out ${
                showPassword ? 'opacity-0' : 'opacity-100'
              }`}
              style={{ transitionProperty: 'opacity' }}
            />
            {/* Password Input (text type) */}
            <Input
              type="text"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className={`${inputClasses} pr-10 w-full absolute inset-0 transition-opacity duration-300 ease-in-out ${
                showPassword ? 'opacity-100' : 'opacity-0'
              }`}
              style={{ transitionProperty: 'opacity' }}
            />
            {/* Toggle Button */}
            <button
              type="button"
              onClick={handlePasswordToggle}
              className="absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500"
              aria-label={showPassword ? 'Hide password' : 'Show password'}
            >
              {showPassword ? <EyeOff className="w-5 h-5" /> : <Eye className="w-5 h-5" />}
            </button>
          </div>
          {/* Login Button */}
          <Button
            type="submit"
            className="w-full h-12 text-md"
            disabled={!isFormValid}
          >
            Login
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}