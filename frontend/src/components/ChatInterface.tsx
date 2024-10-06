'use client'
import React, { useState, useEffect, useRef } from 'react';
import { Card, CardHeader, CardTitle, CardContent } from '@/components/ui/card';
import { Button } from '@/components/ui/button';
import { Textarea } from '@/components/ui/textarea';
import { ScrollArea } from '@/components/ui/scroll-area';
import ThemeToggle from '@/components/ThemeToggle';
import Link from 'next/link';
import { ArrowLeft } from 'lucide-react';

interface Message {
  content: string;
  isUser: boolean;
}

const ChatInterface: React.FC = () => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState('');
  const [isPinned, setIsPinned] = useState(false);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    scrollToBottom();
  }, [messages]);

  const scrollToBottom = () => {
    if (scrollAreaRef.current) {
      const scrollElement = scrollAreaRef.current.querySelector('[data-radix-scroll-area-viewport]');
      if (scrollElement) {
        scrollElement.scrollTop = scrollElement.scrollHeight;
      }
    }
  };

  const handleSend = () => {
    if (inputMessage.trim()) {
      setMessages([...messages, { content: inputMessage, isUser: true }]);
      setInputMessage('');
      // Simulate AI response
      setTimeout(() => {
        setMessages(prev => [...prev, { content: "Lorem ipsum dolor sit amet, consectetur adipiscing elit. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.", isUser: false }]);
      }, 1000);
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputMessage(e.target.value);
    adjustTextareaHeight();
  };

  const adjustTextareaHeight = () => {
    if (textareaRef.current) {
      textareaRef.current.style.height = 'auto';
      const newHeight = Math.min(textareaRef.current.scrollHeight, 4 * 24);
      textareaRef.current.style.height = `${newHeight}px`;
    }
  };

  const bannerContent = (
    <div className={`bg-gray-800 dark:bg-gray-900 text-gray-100 p-2 max-h-[6rem] overflow-y-auto
      ${isPinned ? 'rounded-b-lg' : 'rounded-lg'}`}>
      <p className="text-sm">Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam auctor, nisl nec ultricies ultricies. Sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.</p>
    </div>
  );

  return (
    <div className="w-full h-[calc(100vh-2rem)] flex flex-col">
      <Card className="w-full flex-grow flex flex-col overflow-hidden">
        <CardHeader className={`${isPinned ? 'sticky top-0 z-10' : ''} bg-background flex flex-row items-center justify-between`}>
          <div className="flex items-center">
            <Link href="/dashboard" className="mr-4">
              <Button variant="outline" size="icon">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
            <CardTitle>Accelerator Agent</CardTitle>
          </div>
          <div className="flex items-center space-x-2">
            <Button variant="outline" size="sm" onClick={() => setIsPinned(!isPinned)}>
              {isPinned ? 'Unpin' : 'Pin'}
            </Button>
            <ThemeToggle />
          </div>
        </CardHeader>
        {isPinned && (
          <div className="sticky top-0 left-0 right-0 z-10">
            {bannerContent}
          </div>
        )}
        <CardContent className="flex-grow overflow-hidden">
          <ScrollArea 
            className={`h-full pr-4 ${isPinned ? 'pt-4' : ''}`} 
            ref={scrollAreaRef}
          >
            {!isPinned && bannerContent}
            <div className="flex flex-col space-y-4 pt-6">
              {messages.map((message, index) => (
                <div key={index} className={`flex ${message.isUser ? 'justify-end' : 'justify-start'}`}>
                  <div className={`max-w-[75%] p-5 ${message.isUser ? 'bg-primary text-primary-foreground rounded-tl-3xl rounded-bl-3xl rounded-tr-3xl rounded-br-sm' : 'bg-secondary rounded-tl-3xl rounded-br-3xl rounded-tr-3xl rounded-bl-sm'}`}>
                    <p className="whitespace-pre-wrap break-words">{message.content}</p>
                  </div>
                </div>
              ))}
            </div>
          </ScrollArea>
        </CardContent>
        <div className="p-4 bg-background">
          <div className="flex space-x-2">
            <Textarea
              ref={textareaRef}
              value={inputMessage}
              onChange={handleInputChange}
              placeholder="Type your message..."
              onKeyDown={(e) => e.key === 'Enter' && !e.shiftKey && (e.preventDefault(), handleSend())}
              className="resize-none"
              rows={1}
            />
            <Button onClick={handleSend}>Send</Button>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default ChatInterface;