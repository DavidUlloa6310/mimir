"use client";
import React, {
  useState,
  useEffect,
  useRef,
  useMemo,
  useCallback,
} from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Textarea } from "@/components/ui/textarea";
import { ScrollArea } from "@/components/ui/scroll-area";
import { Separator } from "@/components/ui/separator";
import ThemeToggle from "@/components/ThemeToggle";
import Link from "next/link";
import { ArrowLeft, Send, Shell } from "lucide-react";
import { motion, AnimatePresence, useAnimationControls } from "framer-motion";

interface Message {
  content: string;
  role: "user" | "assistant";
}

const TypewriterText = ({
  text,
  setIsAnimating,
  onCharacterTyped,
}: {
  text: string;
  setIsAnimating: React.Dispatch<React.SetStateAction<boolean>>;
  onCharacterTyped: () => void;
}) => {
  const controls = useAnimationControls();
  const [displayedText, setDisplayedText] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    setIsAnimating(true);
    let currentIndex = 0;
    const intervalDuration = 15;
    const charsPerInterval = 1;

    const interval = setInterval(() => {
      if (currentIndex < text.length) {
        const nextIndex = Math.min(currentIndex + charsPerInterval, text.length);
        setDisplayedText(text.slice(0, nextIndex));
        currentIndex = nextIndex;
        if (containerRef.current) {
          controls.start({ height: containerRef.current.scrollHeight });
        }
        onCharacterTyped(); // Call this after each character is typed
      } else {
        clearInterval(interval);
        setIsAnimating(false);
      }
    }, intervalDuration);

    return () => {
      clearInterval(interval);
      setIsAnimating(false);
    };
  }, [text, controls, setIsAnimating, onCharacterTyped]);

  return (
    <motion.div
      ref={containerRef}
      initial={{ height: 0 }}
      animate={controls}
      transition={{ type: "spring", stiffness: 500, damping: 30 }} // Adjusted for snappier animation
    >
      {displayedText}
    </motion.div>
  );
};

const ChatMessage = ({
  message,
  isNew,
  setIsAnimating,
  onCharacterTyped,
}: {
  message: Message;
  isNew: boolean;
  setIsAnimating: React.Dispatch<React.SetStateAction<boolean>>;
  onCharacterTyped: () => void;
}) => {
  const isAgent = message.role === "assistant";
  const [shouldAnimate, setShouldAnimate] = useState(isNew);

  useEffect(() => {
    if (isNew) {
      setShouldAnimate(true);
    }
  }, [isNew]);

  return (
    <AnimatePresence>
      {shouldAnimate ? (
        <motion.div
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, y: -20 }}
          transition={{ duration: 0.3 }}
          className={`flex ${isAgent ? "justify-start" : "justify-end"} mb-4`}
        >
          <motion.div
            className={`rounded-lg p-4 max-w-[80%] ${
              isAgent ? "bg-primary text-primary-foreground" : "bg-muted"
            }`}
            layout
          >
            <div className="pb-2"> {/* Add padding at the bottom */}
              {isAgent ? (
                <TypewriterText
                  text={message.content}
                  setIsAnimating={setIsAnimating}
                  onCharacterTyped={onCharacterTyped}
                />
              ) : (
                message.content
              )}
            </div>
          </motion.div>
        </motion.div>
      ) : (
        <div
          className={`flex ${isAgent ? "justify-start" : "justify-end"} mb-4`}
        >
          <div
            className={`rounded-lg p-4 max-w-[80%] ${
              isAgent ? "bg-primary text-primary-foreground" : "bg-muted"
            }`}
          >
            <div className="pb-2"> 
              {message.content}
            </div>
          </div>
        </div>
      )}
    </AnimatePresence>
  );
};

interface ChatInterfaceProps {
  threadId: string;
  acceleratorId: string;
}

const ChatInterface: React.FC<ChatInterfaceProps> = ({
  threadId,
  acceleratorId,
}) => {
  const [messages, setMessages] = useState<Message[]>([]);
  const [inputMessage, setInputMessage] = useState("");
  const [isWaiting, setIsWaiting] = useState(false);
  const [newMessageIndex, setNewMessageIndex] = useState(-1);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const [isPinned, setIsPinned] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const [threadTitle, setThreadTitle] = useState("Accelerator Agent");
  const [isAnimating, setIsAnimating] = useState(false);

  const USERNAME = "admin";
  const PASSWORD = "r8RGnqYX=%m0";

  const messagesRef = useRef<Message[]>(messages);

  useEffect(() => {
    messagesRef.current = messages;
  }, [messages]);

  useEffect(() => {
    const fetchThread = async () => {
      try {
        const response = await fetch("http://localhost:8080/chat", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Basic " + btoa(`${USERNAME}:${PASSWORD}`),
          },
          body: JSON.stringify({
            instanceId: "dev274800",
            threadId: threadId,
          }),
        });

        if (response.ok) {
          const data = await response.json();
          setMessages(data.messages);
          if (data.title) {
            setThreadTitle(data.title);
          }
        } else {
          console.error("Failed to fetch thread");
        }
      } catch (error) {
        console.error("Error fetching thread:", error);
      }
    };

    fetchThread();
  }, [threadId]);

  const scrollToBottom = useCallback(() => {
    if (scrollAreaRef.current) {
      const scrollElement = scrollAreaRef.current.querySelector(
        "[data-radix-scroll-area-viewport]"
      );
      if (scrollElement) {
        scrollElement.scrollTop = scrollElement.scrollHeight;
      }
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const handleSend = async () => {
    if (inputMessage.trim() && !isWaiting && !isAnimating) {
      const userMessage = {
        content: inputMessage,
        role: "user",
      };
      const newMessages = [...messages, userMessage];
      setMessages(newMessages);
      setNewMessageIndex(newMessages.length - 1);
      setInputMessage("");
      setIsWaiting(true);
      scrollToBottom();

      try {
        const response = await fetch("http://localhost:8080/chat", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Basic " + btoa(`${USERNAME}:${PASSWORD}`),
          },
          body: JSON.stringify({
            instanceId: "dev274800",
            threadId: threadId,
            message: {
              content: userMessage.content,
            },
            acceleratorId: acceleratorId,
          }),
        });

        if (!response.ok) {
          console.error("Failed to send message");
          setIsWaiting(false);
        }
      } catch (error) {
        console.error("Error sending message:", error);
        setIsWaiting(false);
      }
    }
  };

  const handleInputChange = (e: React.ChangeEvent<HTMLTextAreaElement>) => {
    setInputMessage(e.target.value);
    adjustTextareaHeight();
  };

  const adjustTextareaHeight = () => {
    if (textareaRef.current) {
      textareaRef.current.style.height = "auto";
      const newHeight = Math.min(textareaRef.current.scrollHeight, 4 * 24);
      textareaRef.current.style.height = `${newHeight}px`;
    }
  };

  const bannerContent = (
    <>
      {isPinned && <Separator className="h-1 bg-black dark:bg-gray-700" />}
      <div
        className={`bg-gray-800 dark:bg-gray-900 text-gray-100 p-2 max-h-[6rem] overflow-y-auto
        ${isPinned ? "rounded-b-lg" : "rounded-lg"}`}
      >
        <p className="text-sm">
          Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam
          auctor, nisl nec ultricies ultricies. Sed do eiusmod tempor incididunt
          ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis
          nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo
          consequat.
        </p>
      </div>
    </>
  );

  const memoizedConversation = useMemo(
    () => (
      <div className="flex flex-col space-y-4 pt-6 pb-6">
        {messages.map((message, index) => (
          <ChatMessage
            key={index}
            message={message}
            isNew={index === newMessageIndex}
            setIsAnimating={setIsAnimating}
            onCharacterTyped={scrollToBottom}
          />
        ))}
      </div>
    ),
    [messages, newMessageIndex, setIsAnimating, scrollToBottom]
  );

  useEffect(() => {
    let intervalId: NodeJS.Timeout;

    const pollThread = async () => {
      try {
        const response = await fetch("http://localhost:8080/chat", {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Basic " + btoa(`${USERNAME}:${PASSWORD}`),
          },
          body: JSON.stringify({
            instanceId: "dev274800",
            threadId: threadId,
          }),
        });

        if (response.ok) {
          const data = await response.json();

          // Use messagesRef.current to get the latest messages
          const lastMessage = data.messages[data.messages.length - 1];
          const currentMessages = messagesRef.current;

          if (
            lastMessage.role === "assistant" &&
            (currentMessages.length === 0 ||
              lastMessage.content !== currentMessages[currentMessages.length - 1].content)
          ) {
            setMessages(data.messages);
            setNewMessageIndex(data.messages.length - 1);
            setIsWaiting(false);
            clearInterval(intervalId);
          }
        } else {
          console.error("Failed to fetch thread during polling");
        }
      } catch (error) {
        console.error("Error polling thread:", error);
      }
    };

    if (isWaiting) {
      intervalId = setInterval(pollThread, 3000);
    }

    return () => {
      if (intervalId) {
        clearInterval(intervalId);
      }
    };
  }, [isWaiting, threadId]);

  return (
    <div className="w-full h-[calc(100vh-2rem)] flex flex-col ">
      <Card className="w-full flex-grow flex flex-col overflow-hidden">
        <CardHeader
          className={`${
            isPinned ? "sticky top-0 z-10" : ""
          } bg-background flex flex-row items-center justify-between`}
        >
          <div className="flex items-center">
            <Link href="/dashboard" className="mr-4">
              <Button variant="outline" size="icon">
                <ArrowLeft className="h-4 w-4" />
              </Button>
            </Link>
            <CardTitle>{threadTitle}</CardTitle>
          </div>
          <div className="flex items-center space-x-2">
            <Button
              variant="outline"
              size="sm"
              onClick={() => setIsPinned(!isPinned)}
            >
              {isPinned ? "Unpin" : "Pin"}
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
            className={`h-full pr-4 ${isPinned ? "pt-4" : ""}`}
            ref={scrollAreaRef}
          >
            {!isPinned && bannerContent}
            {memoizedConversation}
            <div className="h-20"></div>
          </ScrollArea>
        </CardContent>
        <div className="p-4 bg-background">
          <div className="flex space-x-2">
            <div className="relative flex-grow">
              <Textarea
                ref={textareaRef}
                value={inputMessage}
                onChange={handleInputChange}
                placeholder={
                  isWaiting || isAnimating
                    ? "Please wait for the agent to respond..."
                    : "Type your message..."
                }
                onKeyDown={(e) =>
                  e.key === "Enter" &&
                  !e.shiftKey &&
                  !(isWaiting || isAnimating) &&
                  (e.preventDefault(), handleSend())
                }
                className="resize-none pr-12"
                rows={1}
                disabled={isWaiting || isAnimating}
              />
              <Button
                onClick={handleSend}
                disabled={isWaiting || isAnimating}
                className="absolute right-4 top-1/2 -translate-y-1/2 h-8 w-8 p-0"
              >
                {(isWaiting || isAnimating) ? (
                  <Shell
                    className="h-4 w-4 animate-spin"
                    style={{ animationDirection: "reverse" }}
                  />
                ) : (
                  <Send className="h-4 w-4" />
                )}
              </Button>
            </div>
          </div>
        </div>
      </Card>
    </div>
  );
};

export default ChatInterface;