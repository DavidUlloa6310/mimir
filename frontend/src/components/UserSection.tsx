"use client";

import { useState, useEffect, useRef } from 'react'
import { useRouter } from 'next/navigation'
import { Avatar, AvatarFallback, AvatarImage } from '@/components/ui/avatar'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuLabel,
  DropdownMenuSeparator,
  DropdownMenuTrigger,
} from "@/components/ui/dropdown-menu";
import { User, LogOut, Upload } from "lucide-react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";

export default function UserSection() {
  const router = useRouter()

  const [avatarSrc, setAvatarSrc] = useState<string | null>(null)
  const [username, setUsername] = useState<string>('')
  const [instanceId, setInstanceId] = useState<string>('')
  const fileInputRef = useRef<HTMLInputElement>(null)

  const handleAvatarUpload = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      const reader = new FileReader();
      reader.onload = (e) => {
        const result = e.target?.result as string;
        setAvatarSrc(result);
        // Save to local storage (as a temporary cache)
        localStorage.setItem("avatarImage", result);
      };
      reader.readAsDataURL(file);
    }
  };

  const openFileExplorer = () => {
    fileInputRef.current?.click();
  };

  const resetAvatar = () => {
    setAvatarSrc(undefined);
    localStorage.removeItem("avatarImage");
  };

  const handleSignOut = () => {
    localStorage.clear()
    router.push('/')
  }

  useEffect(() => {
    const savedAvatar = localStorage.getItem('avatarImage')
    if (savedAvatar) {
      setAvatarSrc(savedAvatar);
    }
    
    const savedUsername = localStorage.getItem('username')
    if (savedUsername) {
      setUsername(savedUsername)
    }

    const savedInstanceId = localStorage.getItem('instanceId')
    if (savedInstanceId) {
      setInstanceId(savedInstanceId)
    }
  }, [])

  return (
    <div className="p-4 h-32 sticky top-0">
      <div className="flex flex-col items-center space-y-4">
        <Avatar className="h-20 w-20">
          <AvatarImage src={avatarSrc || undefined} alt="User avatar" />
          <AvatarFallback>
            <User className="w-16 h-16" />
          </AvatarFallback>
        </Avatar>
        <DropdownMenu>
          <DropdownMenuTrigger asChild>
            <Button variant="outline">{username || 'User'}</Button>
          </DropdownMenuTrigger>
          <DropdownMenuContent className="w-56">
            <DropdownMenuLabel>My Account</DropdownMenuLabel>
            <DropdownMenuSeparator />
            <DropdownMenuItem>
              <User className="mr-2 h-4 w-4" />
              <span>Instance ID: {instanceId || 'N/A'}</span>
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
                    <AvatarImage
                      src={avatarSrc || undefined}
                      alt="User avatar"
                    />
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
              <Button
                variant="destructive"
                className="w-full justify-start"
                onClick={handleSignOut}
              >
                <LogOut className="mr-2 h-4 w-4" />
                <span>Sign out</span>
              </Button>
            </DropdownMenuItem>
          </DropdownMenuContent>
        </DropdownMenu>
      </div>
    </div>
  );
}
