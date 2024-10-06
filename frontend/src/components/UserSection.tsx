'use client'

import { useState, useEffect, useRef } from 'react'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { User, LogOut, Upload } from 'lucide-react'
import { Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger } from '@/components/ui/dialog'

export default function UserSection() {
  const [avatarSrc, setAvatarSrc] = useState<string | null>(null)
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleAvatarUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0]
    if (file) {
      const reader = new FileReader()
      reader.onload = (e) => {
        const result = e.target?.result as string
        setAvatarSrc(result)
        // Save to local storage (as a temporary cache)
        localStorage.setItem('avatarImage', result)
      }
      reader.readAsDataURL(file)
    }
  }

  const openFileExplorer = () => {
    fileInputRef.current?.click()
  }

  const resetAvatar = () => {
    setAvatarSrc(null)
    localStorage.removeItem('avatarImage')
  }

  useEffect(() => {
    // Load avatar from local storage on component mount
    const savedAvatar = localStorage.getItem('avatarImage')
    if (savedAvatar) {
      setAvatarSrc(savedAvatar)
    }
  }, [])

  return (
    <div className="p-4 h-32 sticky top-0">
      <div className="flex flex-col items-center space-y-4">
        <Avatar className="h-20 w-20">
          <AvatarImage src={avatarSrc} alt="User avatar" />
          <AvatarFallback>
            <User className="w-16 h-16" />
          </AvatarFallback>
        </Avatar>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">John Doe</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <User className="mr-2 h-4 w-4" />
              <span>Instance ID: SN001</span>
            </DropdownMenuItem>
            <Dialog>
              <DialogTrigger asChild>
                <DropdownMenuItem onSelect={(event) => event.preventDefault()}>
                  <Upload className="mr-2 h-4 w-4" />
                  <span>Change Avatar</span>
                </DropdownMenuItem>
              </DialogTrigger>
              <DialogContent>
                <DialogHeader>
                  <DialogTitle>Change Avatar</DialogTitle>
                </DialogHeader>
                <div className="flex flex-col items-center space-y-4">
                  <Avatar className="w-32 h-32">
                    <AvatarImage src={avatarSrc || undefined} alt="User avatar" />
                    <AvatarFallback>
                      <User className="w-16 h-16" />
                    </AvatarFallback>
                  </Avatar>
                  <input
                    type="file"
                    accept="image/*"
                    onChange={handleAvatarUpload}
                    className="hidden"
                    ref={fileInputRef}
                  />
                  <Button variant="outline" onClick={openFileExplorer}>
                    Upload New Avatar
                  </Button>
                  <Button variant="destructive" onClick={resetAvatar}>
                    Reset Avatar
                  </Button>
                </div>
              </DialogContent>
            </Dialog>
            <DropdownMenuSeparator />
            <DropdownMenuItem className="flex justify-center">
            <Button variant="destructive" className="w-full justify-start">
              <LogOut className="mr-2 h-4 w-4" />
              <span>Sign out</span>
            </Button>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  )
}