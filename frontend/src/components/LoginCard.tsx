"use client";

import { useState, RefObject } from "react";
import { useRouter } from "next/navigation";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { useToast } from "@/components/hooks/use-toast";
import { motion, AnimatePresence } from "framer-motion";

interface LoginCardProps {
  eyeRef: RefObject<SVGCircleElement>;
  pupilRef: RefObject<SVGCircleElement>;
}

export default function LoginCard({ eyeRef, pupilRef }: LoginCardProps) {
  const router = useRouter();
  const { toast } = useToast();

  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [instanceId, setInstanceId] = useState("");
  const [isPasswordVisible, setIsPasswordVisible] = useState(false);

  const togglePasswordVisibility = () => {
    setIsPasswordVisible(!isPasswordVisible);
  };

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    try {
      const response = await fetch(
        `${process.env.NEXT_PUBLIC_BACKEND_IP}/authorization`,
        {
          method: "POST",
          headers: {
            "Content-Type": "application/json",
            Authorization: "Basic " + btoa(`${username}:${password}`),
          },
          body: JSON.stringify({ instanceId }),
        }
      );

      if (response.ok) {
        localStorage.setItem("instanceId", instanceId);
        localStorage.setItem("username", username);
        localStorage.setItem("password", password);
        router.push("/dashboard");
      } else {
        const errorText = await response.text();
        throw new Error(errorText);
      }
    } catch (error) {
      toast({
        title: "Login Error",
        description: error instanceof Error ? error.message : "An unknown error occurred",
        variant: "destructive",
      });
    }
  };


  return (
    <Card className="w-full backdrop-blur-md bg-white/30 dark:bg-gray-800/90 shadow-xl">
      <CardHeader>
        <CardTitle className="text-2xl font-bold text-center">Login</CardTitle>
        <CardDescription className="text-center">
          Enter your ServiceNow credentials
        </CardDescription>
      </CardHeader>
      <CardContent>
        <form className="space-y-4" onSubmit={handleLogin}>
          <Input
            type="text"
            placeholder="Instance ID"
            value={instanceId}
            onChange={(e) => setInstanceId(e.target.value)}
            className="h-12 bg-white/50 dark:bg-gray-700/50 w-full"
          />
          <div className="border-b dark:bg-slate-800/80 border-gray-700 dark:border-gray-200"></div>
          <Input
            type="text"
            placeholder="Username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="h-12 bg-white/50 dark:bg-gray-700/50 w-full"
          />
          <div className="relative">
            <AnimatePresence mode="wait" initial={false}>
              <motion.div
                key={isPasswordVisible ? "text" : "password"}
                initial={{ filter: "blur(4px)" }}
                animate={{ filter: "blur(0px)" }}
                exit={{ filter: "blur(4px)" }}
                transition={{ duration: 0.2 }}
              >
                <Input
                  type={isPasswordVisible ? "text" : "password"}
                  placeholder="Password"
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                  className="h-12 bg-white/50 dark:bg-gray-700/50 w-full pr-10"
                />
              </motion.div>
            </AnimatePresence>
            <button
              type="button"
              onClick={togglePasswordVisibility}
              className={`absolute inset-y-0 right-0 pr-3 flex items-center text-gray-500 focus:outline-none ${
                !password ? 'cursor-not-allowed opacity-50' : ''
              }`}
              aria-label={isPasswordVisible ? "Hide password" : "Show password"}
              disabled={!password}
            >
              <svg
                viewBox="0 0 24 24"
                fill="none"
                xmlns="http://www.w3.org/2000/svg"
                className="w-5 h-5"
                aria-hidden="true"
              >
                <path
                  className="transition-all duration-200 ease-in-out"
                  d={
                    isPasswordVisible
                      ? "M1 12C1 12 5 4 12 4C19 4 23 12 23 12"
                      : "M1 12C1 12 5 20 12 20C19 20 23 12 23 12"
                  }
                  stroke="currentColor"
                  strokeWidth="1.5"
                  strokeLinecap="round"
                  strokeLinejoin="round"
                />
                <g
                  className={`transition-transform duration-200 ease-in-out ${
                    isPasswordVisible ? "translate-y-0" : "translate-y-1"
                  }`}
                >
                  <circle ref={eyeRef} cy="12" cx="12" r="4" fill="currentColor" />
                  <circle
                    ref={pupilRef}
                    cy="11"
                    cx="13"
                    r="1"
                    fill={isPasswordVisible ? "white" : "currentColor"}
                  />
                  {!isPasswordVisible && (
                    <line
                      x1="19"
                      y1="19"
                      x2="4"
                      y2="4"
                      stroke="currentColor"
                      strokeWidth="2.5"
                      strokeLinecap="round"
                    />
                  )}
                </g>
              </svg>
            </button>
          </div>
          <Button
            type="submit"
            className="w-full h-12 text-md"
            disabled={!username || !password || !instanceId}
          >
            Login
          </Button>
        </form>
      </CardContent>
    </Card>
  );
}
