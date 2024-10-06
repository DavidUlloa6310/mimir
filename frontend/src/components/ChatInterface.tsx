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

const TypewriterText = ({ text }: { text: string }) => {
  const controls = useAnimationControls();
  const [displayedText, setDisplayedText] = useState("");
  const containerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    let currentIndex = 0;
    const interval = setInterval(() => {
      if (currentIndex < text.length) {
        setDisplayedText(text.slice(0, currentIndex + 1));
        currentIndex++;
        if (containerRef.current) {
          controls.start({ height: containerRef.current.scrollHeight });
        }
      } else {
        clearInterval(interval);
      }
    }, 15);

    return () => clearInterval(interval);
  }, [text, controls]);

  return (
    <motion.div
      ref={containerRef}
      initial={{ height: 0 }}
      animate={controls}
      transition={{ type: "spring", stiffness: 200, damping: 20 }}
    >
      {displayedText}
    </motion.div>
  );
};

const ChatMessage = ({
  message,
  isNew,
}: {
  message: Message;
  isNew: boolean;
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
            {isAgent ? (
              <TypewriterText text={message.content} />
            ) : (
              message.content
            )}
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
            {message.content}
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
  const [isFetching, setIsFetching] = useState(false);
  const [newMessageIndex, setNewMessageIndex] = useState(-1);
  const scrollAreaRef = useRef<HTMLDivElement>(null);
  const [isPinned, setIsPinned] = useState(false);
  const textareaRef = useRef<HTMLTextAreaElement>(null);
  const fetchIntervalRef = useRef<NodeJS.Timeout | null>(null);

  const USERNAME = "admin";
  const PASSWORD = "r8RGnqYX=%m0";

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
        } else {
          console.error("Failed to fetch thread");
        }
      } catch (error) {
        console.error("Error fetching thread:", error);
      }
    };

    fetchThread();
  }, [threadId, USERNAME, PASSWORD]);

  // Save messages to localStorage whenever they change
  // useEffect(() => {
  //   localStorage.setItem("chatMessages", JSON.stringify(messages));
  // }, [messages]);

  const scrollToBottom = useCallback(() => {
    if (scrollAreaRef.current) {
      const scrollElement = scrollAreaRef.current.querySelector(
        "[data-radix-scroll-area-viewport]"
      );
      if (scrollElement) {
        scrollElement.scrollTo({
          top: scrollElement.scrollHeight,
          behavior: "smooth",
        });
      }
    }
  }, []);

  useEffect(() => {
    scrollToBottom();
  }, [messages, scrollToBottom]);

  const handleSend = async () => {
    if (inputMessage.trim() && !isWaiting && !isFetching) {
      const userMessage = {
        content: inputMessage,
        role: "user",
      };
      const newMessages = [...messages, userMessage];
      setMessages(newMessages as Message[]);
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

        if (response.ok) {
          const data = await response.json();
          setMessages(data.messages);
          setNewMessageIndex(data.messages.length - 1);
          setIsWaiting(false);
          setIsFetching(false);
        } else {
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
          />
        ))}
      </div>
    ),
    [messages, newMessageIndex]
  );

  // Simulated fetching mechanism
  const simulateFetching = useCallback(() => {
    if (messages.length > 0 && messages[messages.length - 1].role === "user") {
      console.log("Simulating fetch...", new Date().toLocaleTimeString());
      setIsFetching(true);
    } else {
      setIsFetching(false);
      if (fetchIntervalRef.current) {
        clearInterval(fetchIntervalRef.current);
        fetchIntervalRef.current = null;
      }
    }
  }, [messages]);

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
          setMessages(data.messages);
          const lastMessage = data.messages[data.messages.length - 1];
          if (lastMessage.role !== "user") {
            clearInterval(intervalId);
            setIsWaiting(false);
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
  }, [isWaiting, threadId, USERNAME, PASSWORD, simulateFetching]);

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
            <CardTitle>Accelerator Agent</CardTitle>
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
                  isWaiting || isFetching
                    ? "Please wait for the agent to respond..."
                    : "Type your message..."
                }
                onKeyDown={(e) =>
                  e.key === "Enter" &&
                  !e.shiftKey &&
                  !isWaiting &&
                  !isFetching &&
                  (e.preventDefault(), handleSend())
                }
                className="resize-none pr-12"
                rows={1}
                disabled={isWaiting || isFetching}
              />
              <Button
                onClick={handleSend}
                disabled={isWaiting || isFetching}
                className="absolute right-4 top-1/2 -translate-y-1/2 h-8 w-8 p-0"
              >
                {isWaiting || isFetching ? (
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
