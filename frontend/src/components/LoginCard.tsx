'use client'

import { useState } from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Input } from '@/components/ui/input'
import { Button } from '@/components/ui/button'

export default function LoginCard() {
  const [username, setUsername] = useState('')
  const [password, setPassword] = useState('')
  const [instanceId, setInstanceId] = useState('')

  const isFormValid = username && password && instanceId

  return (
    <Card className="w-full backdrop-blur-md bg-white/30 dark:bg-gray-800/30 shadow-xl">
      <CardHeader>
        <CardTitle className="text-2xl font-bold text-center">Login</CardTitle>
        <CardDescription className="text-center">Enter your ServiceNow credentials</CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4">
          <div className="space-y-2">
            <Input
              type="text"
              placeholder="Username"
              value={username}
              onChange={(e) => setUsername(e.target.value)}
              className="bg-white/50 dark:bg-gray-700/50"
            />
          </div>
          <div className="space-y-2">
            <Input
              type="password"
              placeholder="Password"
              value={password}
              onChange={(e) => setPassword(e.target.value)}
              className="bg-white/50 dark:bg-gray-700/50"
            />
          </div>
          <div className="space-y-2">
            <Input
              type="text"
              placeholder="Instance ID"
              value={instanceId}
              onChange={(e) => setInstanceId(e.target.value)}
              className="bg-white/50 dark:bg-gray-700/50"
            />
          </div>
          <Button className="w-full" disabled={!isFormValid}>
            Login
          </Button>
        </form>
      </CardContent>
    </Card>
  )
}